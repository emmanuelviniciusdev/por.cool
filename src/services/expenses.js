import firebase from "firebase/app";
import "firebase/firestore";
import moment from 'moment';

// Services
import userService from './user';

// Helpers
import dateAndTimeHelper from '../helpers/dateAndTime';

const expenses = () => firebase.firestore().collection('expenses');

/**
 * Get all expenses.
 * 
 * @param string userUid 
 * @param Date validity 
 */
const getAll = async (userUid, validity = null) => {
    if (!userUid) return;

    try {
        let allExpenses = expenses()
            .where('user', '==', userUid);

        if (validity !== null) {
            const lastAndNextMonthValidity = dateAndTimeHelper.lastAndNextMonth(validity);

            allExpenses = allExpenses
                .orderBy('spendingDate')
                .where('spendingDate', '<', lastAndNextMonthValidity.startOfNextMonth)
                .where('spendingDate', '>', lastAndNextMonthValidity.endOfLastMonth);
        }

        allExpenses = allExpenses.orderBy('created', 'desc');
        allExpenses = await allExpenses.get();

        return allExpenses.docs.map(expense => {
            return {
                ...expense.data(),
                id: expense.id
            };
        });
    } catch (err) {
        // console.log(err);
        throw new Error(err);
    }
};

/**
 * 
 * @param array expensesToInsert 
 */
const insert = async expensesToInsert => {
    try {
        let batch = firebase.firestore().batch();
        expensesToInsert.forEach(expense => {
            delete expense.id;

            if (expense.differenceAmount !== undefined && expense.differenceAmount <= 0)
                delete expense.differenceAmount;

            batch.set(expenses().doc(), expense)
        });
        await batch.commit();
    } catch (err) {
        throw new Error(err);
    }
};

/**
 * Updates multiple expenses.
 * 
 * OBS: It needs expense to have the document ID as "id" in its object
 * to make update possible.
 * {id: expense_document_id, ...expense}
 * 
 * @param array expensesToUpdate
 */
const bulkUpdate = async expensesToUpdate => {
    try {
        const { FieldValue } = firebase.firestore;

        let batch = firebase.firestore().batch();
        expensesToUpdate.forEach(expense => {
            const { id: expenseDocId } = expense;

            delete expense.id;

            if (expense.differenceAmount !== undefined && expense.differenceAmount <= 0)
                expense.differenceAmount = FieldValue.delete();

            batch.update(expenses().doc(expenseDocId), { ...expense });
        });
        await batch.commit();
    } catch (err) {
        throw new Error(err);
    }
};

/**
 * Updates a single expense.
 * 
 * OBS: It needs expense to have the document ID as "id" in its object
 * to make update possible.
 * {id: expense_document_id, ...expense}
 * 
 * @param object expense 
 * @returns boolean | Promise
 */
const update = async expense => {
    try {
        const { FieldValue } = firebase.firestore;
        const expenseId = expense.id;

        expense.updated = new Date();

        // Delete some data
        delete expense.id;
        if (expense.differenceAmount !== undefined && expense.differenceAmount <= 0)
            expense.differenceAmount = FieldValue.delete();


        await expenses().doc(expenseId).update(expense);
        return true;
    } catch (err) {
        throw new Error(err);
    }
}

/**
 * 
 * @param string expenseId 
 */
const remove = async expenseId => {
    try {
        return await expenses().doc(expenseId).delete();
    } catch (err) {
        throw new Error(err);
    }
};

/**
 * This method will update the property 'lookingAtSpendingDate' of the user
 * to the next month using as reference 'currentLookingAtSpendingDate'.
 * Example:
 * currentLookingAtSpendingDate = 01 January
 * lookingAtSpendingDate = 01 February
 * 
 * And also, if 'autoClone' is true, this method will get all the expenses with type 'invoice' or 'savings'
 * that are in validity with 'currentLookingAtSpendingDate' or have 'indeterminateValidity' setted as true.
 * At the end, it will clone all of them to the new current generated spending date.
 * 
 * @param string userUid 
 * @param Date currentLookingAtSpendingDate
 * @returns boolean
 */
const finishCurrentSpendingDate = async (userUid, currentLookingAtSpendingDate, params = {}) => {
    const { autoClone } = params;

    try {
        const nextLookingAtSpendingDate = moment(currentLookingAtSpendingDate).set('date', 1).add(1, 'months').startOf('day').toDate();
        const expensesToClone = await _getExpensesToClone(userUid, currentLookingAtSpendingDate, nextLookingAtSpendingDate);

        if (autoClone === true) {
            let batchCloneExpenses = firebase.firestore().batch();
            expensesToClone.forEach(expense => {
                if (expense.differenceAmount !== undefined)
                    delete expense.differenceAmount;

                batchCloneExpenses.set(expenses().doc(), expense);
            });
            await batchCloneExpenses.commit();
        }

        await userService.update(userUid, { lookingAtSpendingDate: nextLookingAtSpendingDate });

        return true;
    } catch (err) {
        throw new Error(err);
    }
};

/**
 * Returns expenses that contains type 'invoice' or 'savings' or
 * have the property 'indeterminateValidity' setted as true.
 * 
 * @param string userUid 
 * @returns Promise
 */
const _getExpensesToClone = async (userUid, currentLookingAtSpendingDate, nextLookingAtSpendingDate) => {
    try {
        const expensesToClone = await expenses()
            .where('user', '==', userUid)
            .where('type', 'in', ['invoice', 'savings'])
            .where('status', '==', 'paid')
            .where('spendingDate', '==', currentLookingAtSpendingDate);

        const expensesWithValidity = await expensesToClone.where('validity', '>=', nextLookingAtSpendingDate).get();
        const expensesWithNoValidity = await expensesToClone.where('validity', '==', null).where('indeterminateValidity', '==', true).get();

        const mapUpdateExpenseDates = expense => {
            expense.status = "pending";
            expense.created = new Date();
            expense.spendingDate = nextLookingAtSpendingDate;
            return expense;
        };

        return [
            ...expensesWithValidity.docs.map(expense => mapUpdateExpenseDates(expense.data())),
            ...expensesWithNoValidity.docs.map(expense => mapUpdateExpenseDates(expense.data()))
        ];
    } catch (err) {
        throw new Error(err);
    }
};

/**
 * It returns the years and its months based on 'lookingAtSpendingDate' of the user.
 * This information serves to populate 'FilterByDate' component and is useful
 * in many parts of the application.
 * 
 * @param Object userData 
 * @returns array
 */
const getSpendingDatesList = async ({ userUid, lookingAtSpendingDate }) => {
    if (!userUid) return;

    const allExpenses = await getAll(userUid);

    const { months } = dateAndTimeHelper;

    const currentDate = dateAndTimeHelper.extractOnly(lookingAtSpendingDate, ['year', 'month']);
    const currentDateInObj = {
        year: currentDate.year,
        months: months.slice(0, months.indexOf(currentDate.month) + 1)
    };

    if (allExpenses.length === 0) {
        return [currentDateInObj];
    }

    let spendingDatesList = [];
    let trash = [];

    allExpenses.forEach(({ spendingDate }) => {
        spendingDate = dateAndTimeHelper.transformSecondsToDate(spendingDate.seconds);
        const year = dateAndTimeHelper.extractOnly(spendingDate, ['year']).year;

        if (!trash.includes(year)) {
            spendingDatesList.push({ year, months: year === currentDate.year ? months.slice(0, months.indexOf(currentDate.month) + 1) : months });
            trash.push(year);
        }
    });

    if (!trash.includes(currentDate.year)) {
        spendingDatesList.push(currentDateInObj);
        trash.push(currentDate.year);
    }

    return spendingDatesList;
};

/**
 * Deletes all user's expenses
 * 
 * @param string userUid 
 */
const reset = async (userUid) => {
    try {
        const expensesToDelete = await getAll(userUid);
        
        let batch = firebase.firestore().batch();
        expensesToDelete.forEach(({ id }) => batch.delete(expenses().doc(id)));
        await batch.commit();
    } catch (err) {
        throw new Error(err);
    }
};

export default {
    getAll,
    insert,
    update,
    bulkUpdate,
    remove,
    finishCurrentSpendingDate,
    getSpendingDatesList,
    reset
}
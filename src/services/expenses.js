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
        expensesToInsert.forEach(expense => batch.set(expenses().doc(), expense));
        await batch.commit();
    } catch (err) {
        throw new Error(err);
    }
};

/**
 * It needs expense to have the document ID as "id" in its object
 * to make update possible.
 * {id: expense_document_id, ...expense}
 * 
 * @param array expensesToUpdate
 */
const bulkUpdate = async expensesToUpdate => {
    try {
        let batch = firebase.firestore().batch();
        expensesToUpdate.forEach(expense => {
            const { id: expenseDocId } = expense;
            delete expense.id;
            batch.update(expenses().doc(expenseDocId), { ...expense });
        });
        await batch.commit();
    } catch (err) {
        throw new Error(err);
    }
};

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
            expensesToClone.forEach(expense => batchCloneExpenses.set(expenses().doc(), expense));
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

export default {
    getAll,
    insert,
    bulkUpdate,
    remove,
    finishCurrentSpendingDate
}
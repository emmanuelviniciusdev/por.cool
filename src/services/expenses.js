import firebase from "firebase/app";
import "firebase/firestore";

const expenses = () => firebase.firestore().collection('expenses');

const getAll = async (userUid, validity = null) => {
    try {
        // TODO: Filter by 'spendingDate' as well
        const allExpenses = await expenses()
            .where('user', '==', userUid)
            .orderBy('spendingDate', 'desc')
            .get();
        return allExpenses.docs.map(expense => ({
            id: expense.id,
            ...expense.data()
        }));
    } catch (err) {
        throw new Error(err);
    }
};

const insert = async expensesToInsert => {
    try {
        let batch = firebase.firestore().batch();
        expensesToInsert.forEach(expense => batch.set(expenses().doc(), expense));
        await batch.commit();
    } catch (err) {
        throw new Error(err);
    }
};

const update = async expenseId => {
    try {
        return await expenses().doc(expenseId).update({ ...expenses });
    } catch (err) {
        throw new Error(err);
    }
};

const remove = async expenseId => {
    try {
        return await expenses().doc(expenseId).delete();
    } catch (err) {
        throw new Error(err);
    }
};

export default {
    getAll,
    insert,
    update,
    remove
}
import firebase from "firebase/app";
import "firebase/firestore";

const expenses = () => firebase.firestore().collection('expenses');

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
    insert,
    update,
    remove
}
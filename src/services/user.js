import firebase from "firebase/app";
import "firebase/firestore";

// Services
import expensesService from './expenses';
import balancesService from './balances';

// Helpers
import dateAndTimeHelper from '../helpers/dateAndTime';

const users = () => firebase.firestore().collection('users');

/**
 * Get all user data
 * 
 * @param string uid 
 * @returns Promise
 */
const get = async uid => {
    try {
        const user = await users().doc(uid).get();
        return user.data();
    } catch (err) {
        throw new Error(err);
    }
};

/**
 * Update user data
 * 
 * @param string uid
 * @param object data 
 * @returns Promise
 */
const update = async (uid, data) => {
    try {
        await users().doc(uid).update(data);
        return true;
    } catch (err) {
        throw new Error(err);
    }
}

/**
 * Deletes all expenses and balances from user
 * and set 'lookingAtSpendingDate' to now.
 * 
 * @param string userUid 
 */
const startOver = async userUid => {
    try {
        await expensesService.reset(userUid);
        await balancesService.reset(userUid);
        await users().doc(userUid).update({ lookingAtSpendingDate: dateAndTimeHelper.startOfMonthAndDay(new Date()) });
    } catch (err) {
        throw new Error(err);
    }
};

export default {
    get,
    update,
    startOver
}
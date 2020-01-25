import firebase from "firebase/app";
import "firebase/firestore";
import moment from 'moment';

// Services
import userService from './user';
import expensesServices from './expenses';

const additionalBalances = () => firebase.firestore().collection('additional_balances');
const balanceHistory = () => firebase.firestore().collection('balance_history');

/**
 * It returns how much user have to spend in the given month, that is,
 * the available current balance.
 * 
 * @param object data 
 */
const calculate = async ({ userUid, spendingDate }) => {
    const { monthlyIncome } = await userService.get(userUid);
    const expenses = await expensesServices.getAll(userUid, spendingDate);

    let remainingBalance = monthlyIncome;

    expenses.forEach(({ amount, differenceAmount }) => {
        if (differenceAmount === undefined) differenceAmount = 0;
        remainingBalance -= amount + differenceAmount;
    });

    return remainingBalance;
};

/**
 * Record a data containing information about the user's balance
 * for the given 'spendingDate'.
 * 
 * @param object data
 */
const recordHistory = async ({ userUid, spendingDate }) => {
    try {
        const subtractedSpendingDate = moment(spendingDate).subtract(1, 'months').toDate();

        const { monthlyIncome } = await userService.get(userUid);
        const balance = await calculate({ userUid, spendingDate });
        const lastMonthBalance = await calculate({ userUid, spendingDate: subtractedSpendingDate });

        await balanceHistory().add({
            monthlyIncome,
            balance,
            lastMonthBalance,
            user: userUid,
            spendingDate,
            created: new Date()
        });
    } catch (err) {
        throw new Error(err);
    }
};

export default {
    calculate,
    recordHistory
}
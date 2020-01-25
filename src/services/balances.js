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
    const lastMonthSpendingDate = moment(spendingDate).subtract(1, 'months').toDate();

    const { monthlyIncome } = await userService.get(userUid);
    const currentExpenses = await expensesServices.getAll(userUid, spendingDate);
    const userBalanceHistory = await getHistoryByDate({ userUid, spendingDate: lastMonthSpendingDate });

    let remainingBalance = monthlyIncome;
    currentExpenses.forEach(({ amount, differenceAmount = 0 }) => remainingBalance -= amount + differenceAmount);

    return remainingBalance + (userBalanceHistory.balance !== undefined ? userBalanceHistory.balance : 0);
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
        const lastMonthBalance = await getHistoryByDate({ userUid, spendingDate: subtractedSpendingDate });

        await balanceHistory().add({
            monthlyIncome,
            balance,
            lastMonthBalance: lastMonthBalance.balance !== undefined ? lastMonthBalance.balance : 0,
            user: userUid,
            spendingDate,
            created: new Date()
        });
    } catch (err) {
        throw new Error(err);
    }
};

/**
 * Get an specific balance history by 'spendingDate'
 * 
 * @param object data
 */
const getHistoryByDate = async ({ userUid, spendingDate }) => {
    try {
        const history = await balanceHistory()
            .where('user', '==', userUid)
            .where('spendingDate', '==', spendingDate)
            .get();

        return history.size > 0 ? history.docs[0].data() : {};
    } catch (err) {
        throw new Error(err);
    }
};

export default {
    calculate,
    recordHistory,
    getHistoryByDate
}
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
    try {
        const lastMonthSpendingDate = moment(spendingDate).subtract(1, 'months').toDate();

        const { monthlyIncome } = await userService.get(userUid);
        const currentExpenses = await expensesServices.getAll(userUid, spendingDate);
        const userBalanceHistory = await getHistoryByDate({ userUid, spendingDate: lastMonthSpendingDate });
        const userAdditionalBalances = await calculateAdditionalBalancesOnly({ userUid, spendingDate });

        const lastMonthBalance = userBalanceHistory.balance !== undefined ? _treatFloatNumber(userBalanceHistory.balance) : 0;

        let remainingBalance = _treatFloatNumber(monthlyIncome + lastMonthBalance + userAdditionalBalances);

        // Calculate remaining balance for this month
        currentExpenses.forEach(({ amount, differenceAmount = 0 }) => remainingBalance -= amount + differenceAmount);

        return remainingBalance;
    } catch (err) {
        throw new Error(err);
    }
};

const calculateAdditionalBalancesOnly = async ({ userUid, spendingDate }) => {
    try {
        const userAdditionalBalances = await getAdditionalBalances({ userUid, spendingDate });

        let balancesSum = 0;
        userAdditionalBalances.forEach(data => balancesSum += data.balance);

        return _treatFloatNumber(balancesSum);
    } catch (err) {
        throw new Error(err);
    }
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
 * This function will convert a number with many decimal places
 * to a number with only 2 decimal places.
 * Example:
 * '-2.7284841053187847e-12' => '-2.73'
 * 
 * @param float number 
 */
const _treatFloatNumber = number => {
    let splittedNumber = number.toString().split(".");

    if (splittedNumber[1] === undefined || splittedNumber[1].length < 3)
        return parseFloat(number);

    splittedNumber[1] = splittedNumber[1].substr(0, 3);
    return parseFloat(Number(splittedNumber.join(".")).toFixed(2));
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

/**
 * Inserts a new additional balance.
 * 
 * @param object data 
 */
const addAdditionalBalance = async data => {
    try {
        const { balance, description, spendingDate, userUid: user } = data;
        await additionalBalances().add({
            balance,
            description,
            spendingDate,
            user,
            created: new Date()
        });
    } catch (err) {
        throw new Error(err);
    }
};

/**
 * Returns a list of additional balances by the given 'spendingDate'.
 * 
 * @param object data 
 */
const getAdditionalBalances = async ({ userUid, spendingDate }) => {
    try {
        const balances = await additionalBalances()
            .where('user', '==', userUid)
            .where('spendingDate', '==', spendingDate)
            .orderBy('created', 'desc')
            .get();

        return balances.docs.map(balance => ({ id: balance.id, ...balance.data() }));
    } catch (err) {
        // console.log(err);
        throw new Error(err);
    }
};

/**
 * Deletes an additional balance by document ID
 * 
 * @param string docId 
 */
const removeAdditionalBalance = async docId => {
    try {
        await additionalBalances().doc(docId).delete();
    } catch (err) {
        throw new Error(err);
    }
};

export default {
    calculate,
    calculateAdditionalBalancesOnly,
    recordHistory,
    getHistoryByDate,
    addAdditionalBalance,
    getAdditionalBalances,
    removeAdditionalBalance
}
// Services
import userService from './user';
import expensesServices from './expenses';

/**
 * It returns how much user have to spend in the month, that is,
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

export default {
    calculate
}
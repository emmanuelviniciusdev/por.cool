import dateAndTimeHelper from '../helpers/dateAndTime';

/**
 * Capitalize first letters of user name
 * 
 * @param string name
 * @returns string
 */
const capitalizeName = name => {
    if (name !== undefined) {
        return name
            .split(" ")
            .map(namePart => {
                if (namePart !== "de" && namePart !== "do" && namePart !== "da")
                    namePart = namePart.charAt(0).toUpperCase() + namePart.slice(1);
                return namePart;
            })
            .join(" ");
    }
};

/**
 * Bind the helper 'dateAndTime.extractOnly' in a filter.
 * But attention: this will return only one extracted item at a time.
 * Unlike the helper, that can return an array with all the wishlist
 * passed to it.
 * 
 * @param Date date
 * @param string item
 * @returns object
 */
const extractFromDateOnly = (date, item) => dateAndTimeHelper.extractOnly(date, [item])[item];

/**
 * It adds 'differenceAmount' to 'amount'.
 * 
 * @param integer amount 
 * @param object expense
 * @returns float 
 */
const sumAmounts = (amount, { differenceAmount }) => {
    if (differenceAmount === undefined)
        differenceAmount = 0;

    return parseFloat(amount) + parseFloat(differenceAmount);
};

export default {
    capitalizeName,
    extractFromDateOnly,
    sumAmounts
}
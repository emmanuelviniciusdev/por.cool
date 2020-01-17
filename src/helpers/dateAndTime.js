import moment from 'moment';

const months = ['janeiro', 'fevereiro', 'março', 'abril', 'maio', 'junho', 'julho', 'agosto', 'setembro', 'outubro', 'novembro', 'dezembro'];

/**
 * Extract from date only what is in wishlist
 * 
 * @param Date date 
 * @param array wishlist 
 * @param object params
 * @returns object
 */
const extractOnly = (date, wishlist, params = {}) => {
    const expectedWishlist = ['year', 'month', 'date', 'hour', 'minute', 'second'];

    wishlist.forEach(option => {
        if (!expectedWishlist.includes(option))
            throw new Error(`'${option}' not found in '${expectedWishlist}'`);
    });

    const { monthAsNumber } = params;
    const extractedDate = {};

    wishlist.forEach(option => {
        let foundDate = moment(date).get(option);

        if (option === 'month' && !monthAsNumber)
            foundDate = months[foundDate];

        extractedDate[option] = foundDate;
    });

    return extractedDate;
}

/**
 * Transform determined seconds to Date format
 * 
 * @param integer seconds
 * @param object params 
 * @returns Date
 */
const transformSecondsToDate = (seconds, params = {}) => {
    const { dontSetMidnight } = params;
    let transformedDate = new Date(seconds * 1000);

    if (!dontSetMidnight)
        transformedDate = moment(transformedDate).startOf('day').toDate();

    return transformedDate;
};

/**
 * Receives a date and returns the last day of the last month
 * and the first day of the next month related to this date.
 * 
 * @param Date date
 * @returns object
 */
const lastAndNextMonth = (date) => ({
    endOfLastMonth: moment(date).subtract(1, 'months').endOf('months').toDate(),
    startOfNextMonth: moment(date).add(1, 'months').startOf('months').toDate()
});

/**
 * Set day to 1 and hour to '00:00:00'.
 * 
 * @param Date date 
 * @returns Date
 */
const startOfMonthAndDay = date => moment(date).set('date', 1).startOf('day').toDate();

export default {
    months,
    extractOnly,
    transformSecondsToDate,
    lastAndNextMonth,
    startOfMonthAndDay
}
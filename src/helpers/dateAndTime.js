import moment from 'moment';

const months = ['janeiro', 'fevereiro', 'março', 'abril', 'maio', 'junho', 'julho', 'agosto', 'setembro', 'outubro', 'novembro', 'dezembro'];

/**
 * Extract from date only what is in wishlist
 * 
 * @param Date date 
 * @param array wishlist 
 * @param object params
 * @return object
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
 * @param object firebaseDate 
 * @return Date
 */
const transformSecondsToDate = seconds => new Date(seconds * 1000);

export default {
    extractOnly,
    transformSecondsToDate
}
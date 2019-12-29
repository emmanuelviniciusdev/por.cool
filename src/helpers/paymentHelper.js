import moment from 'moment';

/**
 * Calculate remaining days until user have to pay to use the system again
 * 
 * @param {object} paymentDate 
 * @return integer
 */
const remainingDays = paymentDate => 32 - moment().diff(moment(new Date(paymentDate.seconds * 1000)), 'days');

export default {
    remainingDays
}
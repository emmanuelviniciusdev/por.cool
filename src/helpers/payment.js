import moment from "moment";
import dateAndTimeService from "./dateAndTime";

/**
 * Calculate remaining days until user have to pay to use the system again
 *
 * @param {object} paymentDate
 * @return integer
 */
const remainingDays = paymentDate => {
  const transformedPaymentDate = dateAndTimeService.transformSecondsToDate(
    paymentDate.seconds
  );
  const passedDays = moment().diff(moment(transformedPaymentDate), "days");
  return 32 - passedDays;
};

export default {
  remainingDays
};

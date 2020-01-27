import firebase from "firebase/app";
import "firebase/firestore";

// Helpers
import paymentHelper from '../helpers/payment';
import dateAndTimeHelper from '../helpers/dateAndTime';

/**
 * Returns information about user's last payment
 * 
 * @param string userUid 
 */
const payments = () => firebase.firestore().collection('payments');

/**
 * Returns information about last user's payment
 * 
 * @param string userUid 
 */
const lastPaymentInfo = async userUid => {
    try {
        const lastPayment = await payments()
            .where("user", "==", userUid)
            .orderBy("paymentDate", "desc")
            .limit(1)
            .get();

        const paymentDate = dateAndTimeHelper.transformSecondsToDate(lastPayment.docs[0].data().paymentDate.seconds);

        return {
            paymentDate,
            remainingDays: paymentHelper.remainingDays(lastPayment.docs[0].data().paymentDate)
        };
    } catch (err) {
        throw new Error(err);
    }
};

export default {
    lastPaymentInfo
};
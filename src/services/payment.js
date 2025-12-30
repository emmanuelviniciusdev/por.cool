import firebase from "firebase/app";
import "firebase/firestore";

// Helpers
import paymentHelper from "../helpers/payment";
import dateAndTimeHelper from "../helpers/dateAndTime";

const payments = () => firebase.firestore().collection("payments");

/**
 * Get all user's payment
 *
 * @param string data
 */
const getAll = async userUid => {
  if (!userUid) return;

  try {
    let allPayments = payments().where("user", "==", userUid);

    allPayments = allPayments.orderBy("paymentDate", "desc");
    allPayments = await allPayments.get();

    return allPayments.docs.map(expense => {
      return {
        ...expense.data(),
        id: expense.id
      };
    });
  } catch (err) {
    // console.log(err);
    throw new Error(err);
  }
};

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

    const paymentDate = dateAndTimeHelper.transformSecondsToDate(
      lastPayment.docs[0].data().paymentDate.seconds
    );

    return {
      paymentDate,
      remainingDays: paymentHelper.remainingDays(
        lastPayment.docs[0].data().paymentDate
      )
    };
  } catch (err) {
    throw new Error(err);
  }
};

/**
 * Deletes all user's payments
 *
 * @param string userUid
 */
const reset = async userUid => {
  try {
    const paymentsToDelete = await getAll(userUid);

    let batch = firebase.firestore().batch();
    paymentsToDelete.forEach(({ id }) => batch.delete(payments().doc(id)));

    await batch.commit();
  } catch (err) {
    console.log(err);
    throw new Error(err);
  }
};

export default {
  getAll,
  lastPaymentInfo,
  reset
};

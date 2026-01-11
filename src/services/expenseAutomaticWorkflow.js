import firebase from "firebase/app";
import "firebase/firestore";

// Helpers
import dateAndTimeHelper from "../helpers/dateAndTime";

const expenseAutomaticWorkflow = () =>
  firebase.firestore().collection("expense_automatic_workflow");

/**
 * Get all automatic workflow entries for a user filtered by spending date (month/year).
 *
 * @param string userUid
 * @param Date spendingDate
 * @returns Promise<Array>
 */
const getAll = async (userUid, spendingDate = null) => {
  if (!userUid) return [];

  try {
    let query = expenseAutomaticWorkflow().where("user", "==", userUid);

    if (spendingDate !== null) {
      const lastAndNextMonthValidity = dateAndTimeHelper.lastAndNextMonth(
        spendingDate
      );

      query = query
        .orderBy("spendingDate")
        .where("spendingDate", "<", lastAndNextMonthValidity.startOfNextMonth)
        .where("spendingDate", ">", lastAndNextMonthValidity.endOfLastMonth);
    }

    query = query.orderBy("created", "desc");
    const snapshot = await query.get();

    return snapshot.docs.map(doc => ({
      ...doc.data(),
      id: doc.id
    }));
  } catch (err) {
    throw new Error(err);
  }
};

/**
 * Insert a new automatic workflow entry.
 *
 * @param object data - { user, description, base64_image, spendingDate, syncStatus }
 * @returns Promise<string> - The document ID
 */
const insert = async data => {
  try {
    const docRef = await expenseAutomaticWorkflow().add({
      ...data,
      created: new Date(),
      syncStatus: "pending",
      extracted_expense_content_from_image: [],
      onPremiseSyncDatetime: null
    });
    return docRef.id;
  } catch (err) {
    throw new Error(err);
  }
};

/**
 * Update an automatic workflow entry.
 *
 * @param string docId
 * @param object data
 * @returns Promise<boolean>
 */
const update = async (docId, data) => {
  try {
    await expenseAutomaticWorkflow()
      .doc(docId)
      .update({
        ...data,
        updated: new Date()
      });
    return true;
  } catch (err) {
    throw new Error(err);
  }
};

/**
 * Delete an automatic workflow entry.
 *
 * @param string docId
 * @returns Promise
 */
const remove = async docId => {
  try {
    return await expenseAutomaticWorkflow()
      .doc(docId)
      .delete();
  } catch (err) {
    throw new Error(err);
  }
};

/**
 * Delete all automatic workflow entries for a user.
 *
 * @param string userUid
 * @returns Promise
 */
const reset = async userUid => {
  try {
    const entries = await getAll(userUid);

    let batch = firebase.firestore().batch();
    entries.forEach(({ id }) =>
      batch.delete(expenseAutomaticWorkflow().doc(id))
    );
    await batch.commit();
  } catch (err) {
    throw new Error(err);
  }
};

export default {
  getAll,
  insert,
  update,
  remove,
  reset
};

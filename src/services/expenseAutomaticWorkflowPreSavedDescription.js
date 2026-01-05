import firebase from "firebase/app";
import "firebase/firestore";

const preSavedDescriptions = () =>
  firebase
    .firestore()
    .collection("expense_automatic_workflow_pre_saved_description");

/**
 * Get all pre-saved descriptions for a user.
 *
 * @param string userUid
 * @returns Promise<Array>
 */
const getAll = async userUid => {
  if (!userUid) return [];

  try {
    const snapshot = await preSavedDescriptions()
      .where("user", "==", userUid)
      .orderBy("created", "desc")
      .get();

    return snapshot.docs.map(doc => ({
      ...doc.data(),
      id: doc.id
    }));
  } catch (err) {
    throw new Error(err);
  }
};

/**
 * Insert a new pre-saved description.
 *
 * @param object data - { user, description }
 * @returns Promise<string> - The document ID
 */
const insert = async data => {
  try {
    const docRef = await preSavedDescriptions().add({
      ...data,
      created: new Date(),
      onPremiseSyncDatetime: null
    });
    return docRef.id;
  } catch (err) {
    throw new Error(err);
  }
};

/**
 * Delete a pre-saved description.
 *
 * @param string docId
 * @returns Promise
 */
const remove = async docId => {
  try {
    return await preSavedDescriptions()
      .doc(docId)
      .delete();
  } catch (err) {
    throw new Error(err);
  }
};

/**
 * Delete all pre-saved descriptions for a user.
 *
 * @param string userUid
 * @returns Promise
 */
const reset = async userUid => {
  try {
    const entries = await getAll(userUid);

    let batch = firebase.firestore().batch();
    entries.forEach(({ id }) => batch.delete(preSavedDescriptions().doc(id)));
    await batch.commit();
  } catch (err) {
    throw new Error(err);
  }
};

export default {
  getAll,
  insert,
  remove,
  reset
};

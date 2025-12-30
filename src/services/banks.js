import firebase from "firebase/app";
import "firebase/firestore";

const banksCollection = () => firebase.firestore().collection("banks");

/**
 * Inserts a new bank/institution.
 *
 * @param object data
 */
const addBank = async data => {
  try {
    const {
      nome,
      cartaoCredito,
      movimentacaoDinheiro,
      investimentos,
      observacoes,
      userUid: user
    } = data;
    await banksCollection().add({
      nome,
      cartaoCredito,
      movimentacaoDinheiro,
      investimentos,
      observacoes,
      user,
      created: new Date()
    });
  } catch (err) {
    console.error("Error adding bank:", err);
    throw err;
  }
};

/**
 * Returns a list of banks/institutions for the given user.
 *
 * @param string userUid
 */
const getBanks = async userUid => {
  try {
    const banks = await banksCollection()
      .where("user", "==", userUid)
      .orderBy("created", "desc")
      .get();

    return banks.docs.map(bank => ({ id: bank.id, ...bank.data() }));
  } catch (err) {
    console.error("Error getting banks:", err);
    throw err;
  }
};

/**
 * Updates a bank/institution by document ID
 *
 * @param object data
 */
const updateBank = async data => {
  try {
    const {
      id,
      nome,
      cartaoCredito,
      movimentacaoDinheiro,
      investimentos,
      observacoes
    } = data;
    await banksCollection()
      .doc(id)
      .update({
        nome,
        cartaoCredito,
        movimentacaoDinheiro,
        investimentos,
        observacoes,
        updated: new Date()
      });
  } catch (err) {
    throw new Error(err);
  }
};

/**
 * Deletes a bank/institution by document ID
 *
 * @param string docId
 */
const removeBank = async docId => {
  try {
    await banksCollection()
      .doc(docId)
      .delete();
  } catch (err) {
    throw new Error(err);
  }
};

export default {
  addBank,
  getBanks,
  updateBank,
  removeBank
};

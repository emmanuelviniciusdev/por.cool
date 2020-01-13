import firebase from "firebase/app";
import "firebase/firestore";

const users = () => firebase.firestore().collection('users');

/**
 * Get all user data
 * 
 * @param string uid 
 * @returns Promise
 */
const get = async uid => {
    try {
        const user = await users().doc(uid).get();
        return user.data();
    } catch (err) {
        throw new Error(err);
    }
};

/**
 * Update user data
 * 
 * @param string uid
 * @param object data 
 * @returns Promise
 */
const update = async (uid, data) => {
    try {
        await users().doc(uid).update(data);
        return true;
    } catch (err) {
        throw new Error(err);
    }
}

export default {
    get,
    update
}
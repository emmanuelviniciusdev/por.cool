import firebase from "firebase/app";
import "firebase/firestore";

const users = () => firebase.firestore().collection('users');

const get = async uid => {
    try {
        const user = await users().doc(uid).get();
        return user.data();
    } catch (err) {
        throw new Error(err);
    }
};

export default {
    get
}
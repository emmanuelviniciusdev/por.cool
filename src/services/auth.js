import firebase from "firebase/app";
import "firebase/firestore";

// Services
import userService from '../services/user';

const auth = () => firebase.auth();

/**
 * Change user's password
 * 
 * @param string currentPassword 
 * @param string newPassword 
 */
const changePassword = async (currentPassword, newPassword) => {
    try {
        const user = auth().currentUser;
        
        try {
            await reauthenticate(currentPassword);
        } catch {
            return _response({ error: true, message: "a senha está incorreta" });
        }

        await user.updatePassword(newPassword);
        return _response({ message: "senha alterada com sucesso" });
    } catch (err) {
        console.log(err);
        return _response({ error: true, message: "ocorreu um erro ao alterar a senha [1]" });
    }
};


/**
 * Change user's email
 * 
 * @param string password 
 * @param string newEmail 
 */
const changeEmail = async (password, newEmail) => {
    try {
        const user = auth().currentUser;

        try {
            await reauthenticate(password);
        } catch {
            return _response({ error: true, message: "a senha está incorreta" });
        }

        if (newEmail === user.email)
            return _response({ error: true, message: "então quer dizer que seu novo e-mail é igual ao atual? 😞" });

        await user.updateEmail(newEmail);
        await userService.update(user.uid, { email: newEmail });

        return _response({ message: "e-mail alterado com sucesso" });
    } catch (err) {
        return _response({ error: true, message: "ocorreu um erro ao alterar o seu e-mail [1]" });
    }
};

/**
 * Re-authenticate user based on user's password
 * 
 * @param string password 
 */
const reauthenticate = async password => {
    try {
        const user = auth().currentUser;
        const credential = firebase.auth.EmailAuthProvider.credential(user.email, password);

        return await user.reauthenticateWithCredential(credential);
    } catch (err) {
        throw err;
    }
};

/**
 * Factory to generate responses
 * 
 * @param object params 
 */
const _response = ({ error, message }) => ({ error: error ?? false, message: message ?? "" });

export default {
    changePassword,
    changeEmail,
    reauthenticate
}
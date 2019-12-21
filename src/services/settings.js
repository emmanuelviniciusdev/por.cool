import firebase from "firebase/app";
import "firebase/database";

const messages = {
    maintenance: "Uma pequena manutenção de rotina está sendo realizada neste exato momento. Por favor, tente novamente mais tarde.",
    blockUserRegistration: "o registro de novos usuários encontra-se temporariamente indisponível",
    unexpectedError: "Ocorreu um erro inesperado. Por favor, tente novamente mais tarde"
};

/**
 * @return object
 */
const checkMaintenances = async (filterBy = "") => {
    let status = {};

    try {
        const settings = await firebase
            .database()
            .ref(`settings/${filterBy}`)
            .once("value");

        if (filterBy !== "") {
            let message = (settings.val() !== null && settings.val()) ? (messages[filterBy] ? messages[filterBy] : messages.unexpectedError) : (false);
            status[filterBy] = message;
            return status;
        }

        Object.keys(settings.val()).map(setting => {
            status[setting] = (settings.val()[setting]) ? (messages[setting] ? messages[setting] : messages.unexpectedError) : false;
        });

        return status;
    } catch (err) {
        throw new Error(err);
    }
}

export default {
    checkMaintenances
}
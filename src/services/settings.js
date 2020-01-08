import firebase from "firebase/app";
import "firebase/firestore";

const messages = {
    maintenance: "Uma pequena manutenção de rotina está sendo realizada neste exato momento. Por favor, tente novamente mais tarde.",
    blockUserRegistration: "o registro de novos usuários encontra-se temporariamente indisponível",
    unexpectedError: "Ocorreu um erro inesperado. Por favor, tente novamente mais tarde"
};

const checkMaintenances = async (filterBy = "") => {
    let status = {};

    try {
        const settingsReq = await firebase
            .firestore()
            .collection(`settings`)
            .limit(1)
            .get();
        const settings = settingsReq.docs[0].data();

        if (filterBy !== "") {
            let message = (settings !== null && settings) ? (messages[filterBy] ? messages[filterBy] : messages.unexpectedError) : (false);
            status[filterBy] = message;
            return status;
        }

        Object.keys(settings).map(setting => {
            status[setting] = (settings[setting]) ? (messages[setting] ? messages[setting] : messages.unexpectedError) : false;
        });

        return status;
    } catch (err) {
        throw new Error(err);
    }
}

export default {
    checkMaintenances
}
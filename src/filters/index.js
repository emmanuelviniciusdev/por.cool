/**
 * Capitalize first letters of user name
 * 
 * @param string name
 * @returns string
 */
const capitalizeName = name => {
    return name
        .split(" ")
        .map(namePart => {
            if (namePart !== "de" && namePart !== "do" && namePart !== "da")
                namePart = namePart.charAt(0).toUpperCase() + namePart.slice(1);
            return namePart;
        })
        .join(" ");
};

export default {
    capitalizeName
}
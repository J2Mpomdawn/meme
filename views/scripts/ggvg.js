/**generate groups process */
function generateGroups(cs) {
    //combo box
    const groups = document.getElementById("group");
    
    //clear all
    for (let i = groups.options.length - 1; i >= 0; i--) {
        groups.remove(i);
    }

    //get groups
    fetch("https://api.mentemori.icu/wgroups")
        .then(resp => resp.json())
        .then((grps) => {
            //make and add choices
            for (const group of grps.data) {
                if (!group.globalgvg) {
                    continue;
                }

                if (!group.worlds.map(x => Math.floor(x / 1000) * 1000 == cs)
                    .reduce((x, y) => x || y)
                    ) {
                    continue;
                }

                //create option
                const opt = document.createElement("option");
                opt.value = group.group_id;
                opt.text = group.group_id + " (" + group.worlds.map(x => "W" + x % 1000).join(", ") + ")";

                groups.add(opt);
            }
        });
}

/**page load completion event*/
document.addEventListener("DOMContentLoaded", () => {
    //get groups
    generateGroups(document.getElementById("server").value * 1000)
});

/**server changed event */
document.getElementById("server").addEventListener("change", (e) => {
    //get groups
    generateGroups(e.target.value * 1000);
});
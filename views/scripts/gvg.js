/**generate worlds process */
function generateWorlds(cs) {
    //combo box
    const worlds = document.getElementById("world");
    
    //clear all
    for (let i = worlds.options.length - 1; i >= 0; i--) {
        worlds.remove(i);
    }

    //current server
    switch (cs*1) {
        case 1:
            cs = "jp"

            break;
        case 2:
            cs = "kr"
            
            break;
        case 3:
            cs = "as"
            
            break;
        case 4:
            cs = "na"
            
            break;
        case 5:
            cs = "en"
            
            break;
        case 6:
            cs = "gl"
            
            break;
        default:
            break;
    }

    //get worlds
    fetch("https://api.mentemori.icu/worlds")
        .then(resp => resp.json())
        .then((wlds) => {
            //make and add choices
            for (const world of wlds.data) {
                if (!world.localgvg) {
                    continue;
                }

                if (world.server != cs) {
                    continue;
                }

                //create option
                const opt = document.createElement("option");
                opt.value = world.world_id % 1000;
                opt.text = "W" + world.world_id % 1000

                worlds.add(opt);
            }
        });
}

/**page load completion event*/
document.addEventListener("DOMContentLoaded", () => {
    //get worlds
    generateWorlds(document.getElementById("server").value);
});

/**server changed event */
document.getElementById("server").addEventListener("change", (e) => {
    //get worlds
    generateWorlds(e.target.value);
});
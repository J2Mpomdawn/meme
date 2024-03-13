/**page load completion event*/
document.addEventListener("DOMContentLoaded", () => {
    //combo box
    const worlds = document.getElementById("world");

    //make and add coices
    for (let i = 1; i <= 111; i++) {
        //create option
        const opt = document.createElement("option");
        opt.value = i;
        opt.text = "W" + i;

        worlds.add(opt);
    }
});
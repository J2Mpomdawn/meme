document.getElementById("b1").addEventListener("click", ()=>{
    post();
});

post = () => {
    const ajax = new XMLHttpRequest;
    ajax.open("POST", "/cmdtest/exec", true);
    ajax.setRequestHeader("content-type", "application/json");
    //document.getElementById("i1").value
    ajax.send(JSON.stringify(
        {
            "app": "echo",
            "args": [
                "%path%",
                "$path"
            ]
        }
    ));
};
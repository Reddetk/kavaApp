import React, { useState } from "react";



const App = () => {
    let [likes, SetLikes] = useState(0);

    function increment(){
        SetLikes(likes++);
        console.log(likes);
    }   
    
    function decrement(){
        SetLikes(likes--);
        console.log(likes);
    }   

    function Save(){
        likesArr.push(likes);
        SetLikes(likes = 0);
        console.log(likes, likesArr);
    }

    let likesArr = [];

    return (
        <div>
            <h1>{likes}</h1>
            <button onClick={increment}>какшки</button>
            <button onClick={decrement}>не в кармашки</button>
            <button onClick={Save}>не в кармашки</button>
        </div>
    );
}


export default App;
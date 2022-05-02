import {Board} from "../api-ts/generated/models";
import {BoardDecoder} from "../api-ts/generated/decoders";

async function getBoard() : Promise<Board> {

    const result = await fetch("/board"); 
    const json = await result.json();
    

    const decoded = BoardDecoder.decode(json); 
    return decoded; 
}
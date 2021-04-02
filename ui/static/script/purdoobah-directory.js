// search bar, attach type input listener
const inputSearch = document.getElementById("inputSearch");
inputSearch.addEventListener("input", search);

// grab all purdoobah cards
const purdoobahCards = document.getElementsByClassName("purdoobah-card");

// the archives-incomplete gif
// set it to invisible initially and pause it
const archivesIncomplete = document.getElementsByClassName("archives-incomplete")[0];
archivesIncomplete.style.display = "none";
archivesIncomplete.pause();

// `search` shows/hides Purdoobah Cards based on the search term
function search() {
    const searchTerm = this.value.toLowerCase();

    // show/hide purdoobah cards depending on whether or not they match the search term
    const cardsHidden = togglePurdoobahCards(searchTerm);

    // if every purdoobah card is hidden, show archives-incomplete gif
    toggleArchivesIncomplete(cardsHidden);
}

function togglePurdoobahCards(searchTerm) {
    let cardsHidden = 0;
    for (let i = 0; i < purdoobahCards.length; i++) {
        const purdoobahCard = purdoobahCards.item(i);

        const name = purdoobahCard.dataset.name;
        const emoji = purdoobahCard.dataset.emoji;
        const birthCertificateName = JSON.parse(purdoobahCard.dataset.birthCertificateName);
        const yearsMarched = purdoobahCard.dataset.yearsMarched;

        if (containsSearch(searchTerm, name, emoji, birthCertificateName, yearsMarched)) {
            if (purdoobahCard.style.display === "block") {
                // purdoobah card is already shown, don't do anything
                continue;
            }

            purdoobahCard.style.display = "block";
        } else {
            if (purdoobahCard.style.display === "none") {
                // purdoobah card is already hidden, don't do anything
                cardsHidden++;
                continue;
            }

            purdoobahCard.style.display = "none";
            cardsHidden++;
        }
    }

    return cardsHidden;
}

// `containsSearch` returns true if search term can be found in any of the provided values
function containsSearch(searchTerm, name, emoji, birthCertificateName, yearsMarched) {
    // name
    if (name.toLowerCase().includes(searchTerm)) {
        return true;
    }

    // emoji
    if (emoji.includes(searchTerm)) {
        return true;
    }

    // first name
    if (birthCertificateName.first.toLowerCase().includes(searchTerm)) {
        return true;
    }

    // middle name
    if (birthCertificateName.hasOwnProperty("middle") &&
        birthCertificateName.middle.toLowerCase().includes(searchTerm)) {
        return true;
    }

    // last name
    if (birthCertificateName.last.toLowerCase().includes(searchTerm)) {
        return true;
    }

    // years marched
    if (yearsMarched.includes(searchTerm)) {
        return true;
    }

    return false;
}

// `toggleArchivesIncomplete` shows/hides the archives-incomplete video
// depending on if the search term returns zero results
function toggleArchivesIncomplete(cardsHidden) {
    if (cardsHidden === purdoobahCards.length) {
        if (archivesIncomplete.style.display === "block") {
            // video is already playing, don't do anything
            return;
        }
        // restart the video from the beginning, show it, then play it
        archivesIncomplete.currentTime = 0;
        archivesIncomplete.style.display = "block";
        archivesIncomplete.play();
    } else {
        if (archivesIncomplete.style.display === "none") {
            // video is already hidden, don't do anything
            return;
        }
        // hide the video, then pause it
        archivesIncomplete.style.display = "none";
        archivesIncomplete.pause();
    }
}

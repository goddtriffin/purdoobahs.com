// search bar, attach type input listener
const inputSearch = document.getElementById("inputSearch");
inputSearch.addEventListener("input", search);

// `search` shows/hides Purdoobah Cards based on the search term
function search() {
    const search = this.value.toLowerCase();

    const purdoobahCards = document.getElementsByClassName("purdoobahCard");

    for (let i = 0; i < purdoobahCards.length; i++) {
        const purdoobahCard = purdoobahCards.item(i);

        const name = purdoobahCard.dataset.name;
        const emoji = purdoobahCard.dataset.emoji;
        const birthCertificateName = JSON.parse(purdoobahCard.dataset.birthCertificateName);
        const yearsMarched = purdoobahCard.dataset.yearsMarched;

        if (containsSearch(search, name, emoji, birthCertificateName, yearsMarched)) {
            purdoobahCard.style.display = "block";
        } else {
            purdoobahCard.style.display = "none";
        }
    }
}

// `containsSearch` returns true if search term can be found in any of the provided values
function containsSearch(search, name, emoji, birthCertificateName, yearsMarched) {
    // name
    if (name.toLowerCase().includes(search)) {
        return true;
    }

    // emoji
    if (emoji.includes(search)) {
        return true;
    }

    // first name
    if (birthCertificateName.first.toLowerCase().includes(search)) {
        return true;
    }

    // middle name
    if (birthCertificateName.hasOwnProperty("middle") &&
        birthCertificateName.middle.toLowerCase().includes(search)) {
        return true;
    }

    // last name
    if (birthCertificateName.last.toLowerCase().includes(search)) {
        return true;
    }

    // years marched
    if (yearsMarched.includes(search)) {
        return true;
    }

    return false;
}

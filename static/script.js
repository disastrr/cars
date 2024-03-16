// Search function
function search() {
    var input = document.getElementById('searchInput').value.toLowerCase();
    var cards = document.getElementsByClassName('card');

    for (var i = 0; i < cards.length; i++) {
        var modelName = cards[i].getElementsByClassName('modelName')[0].innerText.toLowerCase();
        var manufacturerName = cards[i].getElementsByClassName('manufacturerName')[0].innerText.toLowerCase();
        var category = cards[i].getElementsByClassName('category')[0].innerText.toLowerCase();
        var year = cards[i].getElementsByClassName('year')[0].innerText.toLowerCase();

        // Check if input matches any of the criteria
        if (modelName.includes(input) || manufacturerName.includes(input) || category.includes(input) || year.includes(input)) {
            cards[i].style.display = "";
        } else {
            cards[i].style.display = "none";
        }
    }
}

function hideUnselected() {
    var checkboxes = document.querySelectorAll('input[type="checkbox"]');
    var selectedIds = [];
    
    checkboxes.forEach(function(checkbox) {
        if (checkbox.checked) {
            selectedIds.push(checkbox.id);
        }
    });

    var cards = document.querySelectorAll('.card');
    cards.forEach(function(card) {
        if (!selectedIds.includes(card.querySelector('input[type="checkbox"]').id)) {
            card.style.display = 'none';
            card.classList.add('highlight');
        } else {
            card.style.display = 'block';
            card.classList.remove('highlight');
        }
    });
}


// Wrap JavaScript code inside window.onload to ensure it executes after the DOM is loaded
window.onload = function() {
    // Add event listener to the search input field
    document.getElementById('searchInput').addEventListener('input', function() {
        search(); // Call the search function whenever the input value changes
    });

    // Add event listener to the compare button
    document.getElementById('compare').addEventListener('click', hideUnselected);

    // Get all collapsible buttons
    var coll = document.getElementsByClassName("collapsible");

    // Add click event listener to each button
    for (var i = 0; i < coll.length; i++) {
        coll[i].addEventListener("click", function() {
            // Toggle the visibility of the content when the button is clicked
            this.classList.toggle("active");
            var content = this.nextElementSibling;
            if (content.style.display === "block") {
                content.style.display = "none";
            } else {
                content.style.display = "block";
            }
        });
    }

    // popover
    var popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'))
    var popoverList = popoverTriggerList.map(function (popoverTriggerEl) {
        return new bootstrap.Popover(popoverTriggerEl)
    });
};



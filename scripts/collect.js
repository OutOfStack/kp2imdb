// Select all film items on the page
const filmItems = document.querySelectorAll('.item');
const filmsData = [];

// Loop through each film element and extract data
filmItems.forEach(item => {
  // Extract the Russian title from .nameRus element
  const nameElement = item.querySelector('.nameRus');
  let nameText = nameElement ? nameElement.innerText.trim() : '';

  // Remove trailing parentheses content and extract year if applicable
  let year = "";
  const parenMatch = nameText.match(/\s*\(([^)]+)\)$/);
  if (parenMatch) {
    const content = parenMatch[1];
    // If the content is exactly a four-digit year, set it as the year
    const simpleYearMatch = content.match(/^(\d{4})$/);
    if (simpleYearMatch) {
      year = simpleYearMatch[1];
    }
    // Remove the parentheses content from the name regardless
    nameText = nameText.replace(/\s*\([^)]+\)$/, "");
  }

  // Extract the English title from .nameEng element
  const originalNameElement = item.querySelector('.nameEng');
  const originalName = originalNameElement ? originalNameElement.innerText.trim() : '';

  // Extract your rating from the vote widget element (.vote_widget .myVote)
  const myRatingElement = item.querySelector('.vote_widget .myVote');
  const myRating = myRatingElement ? myRatingElement.innerText.trim() : '';

  // Don't add entry if both name and original_name are empty
  if (!nameText && !originalName) {
    return;
  }

  // Create an object for this film's data
  filmsData.push({
    name: nameText,
    original_name: originalName,
    my_rating: myRating,
    year: year
  });
});

// Output the resulting JSON data in the console
console.log(JSON.stringify(filmsData, null, 2));

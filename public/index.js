document.addEventListener("DOMContentLoaded", function() {
  let dt = null;
    fetch('data.json')
    .then(response => response.json())
    .then(data => {
  dt = data;
      })
    .catch(error => console.error('Error fetching JSON:', error));
});


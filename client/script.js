async function postData(url = '', data = {}) {
    return fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
    })
    .catch(error => {
        console.error('Error:', error.message);
        throw error;
    });
}

document.querySelector('form').addEventListener('submit', function(event) {

    event.preventDefault();

    const url = 'http://localhost:8082/url/'; 
    const data = { url: document.getElementById('linkInput').value }; 

    postData(url, data)
        .then(responseData => {
            const aliasValue = responseData.alias; 
            const resUrl = 'http://localhost:8082/url/' + aliasValue;
            document.getElementById('result').innerHTML = '<a href="' + resUrl + '">' + resUrl + '</a>';
        })
        .catch(error => {
            console.error('Error in postData:', error.message); 
        });
});

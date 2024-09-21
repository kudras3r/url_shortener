
function postData(url = '', data = {}) {
    try {
        const response = fetch(url, {
            mode: "no-cors",
            method: 'POST', 
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        if (!response.ok) {
            throw new Error(`Ошибка ${response.status}: ${response.statusText}`); 
        }

        const result = response.json(); 
        return result; 
    } catch (error) {
        console.error('Ошибка при отправке данных:', error);
        return { error: error.message }; 
    }
}

document.getElementById('linkForm').addEventListener('submit', function(event) {
    event.preventDefault(); 
    const link = document.getElementById('linkInput').value; 

    const apiUrl = 'http://localhost:8082/url/'; 
    
    const result = postData(apiUrl, { "url": link }); 
    console.log(result)
    document.getElementById('result').innerText = "http://localhost:8082/" + result.alias; 
});

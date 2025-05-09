async function getWeather() {
    const latitude = document.getElementById('latitude').value;
    const longitude = document.getElementById('longitude').value;
    const errorDiv = document.getElementById('error');
    const weatherResult = document.getElementById('weather-result');

    if (!latitude || !longitude) {
        showError('Please enter both latitude and longitude');
        return;
    }

    try {
        const response = await fetch(`http://localhost:8080/weather?lat=${latitude}&lon=${longitude}`);
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to fetch weather data');
        }

        displayWeather(data);
        errorDiv.classList.add('hidden');
    } catch (error) {
        showError(error.message);
        weatherResult.classList.add('hidden');
    }
}

function displayWeather(data) {
    const weatherResult = document.getElementById('weather-result');
    const city = document.getElementById('city');
    const temperature = document.getElementById('temperature');
    const description = document.getElementById('description');
    const humidity = document.getElementById('humidity');
    const windSpeed = document.getElementById('windSpeed');

    city.textContent = data.city;
    temperature.textContent = `${data.temperature}°C`;
    description.textContent = data.description;
    humidity.textContent = `${data.humidity}%`;
    windSpeed.textContent = `${data.windSpeed} m/s`;

    weatherResult.classList.remove('hidden');
}

function showError(message) {
    const errorDiv = document.getElementById('error');
    errorDiv.textContent = message;
    errorDiv.classList.remove('hidden');
} 
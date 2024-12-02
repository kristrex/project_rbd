document.addEventListener('DOMContentLoaded', () => {
    const baseURL = 'http://localhost:8080';
    const token = localStorage.getItem('token');

    if (token) {
        try {
            const decodedToken = JSON.parse(atob(token.split('.')[1]));
            const currentTime = Math.floor(Date.now() / 1000);
            if (decodedToken.exp < currentTime) {
                localStorage.removeItem('token');
                showAuthForms();
            } else {
                showContent();
                loadData();
            }
        } catch (error) {
            console.error('Error decoding token:', error);
            localStorage.removeItem('token');
            showAuthForms();
        }
    } else {
        showAuthForms();
    }

    function showAuthForms() {
        document.getElementById('auth-forms').style.display = 'block';
        document.getElementById('content').style.display = 'none';
    }

    function showContent() {
        document.getElementById('auth-forms').style.display = 'none';
        document.getElementById('content').style.display = 'block';
    }

    async function fetchData(endpoint) {
        console.log(`Fetching data from ${endpoint}`);
        const response = await fetch(`${baseURL}/api/${endpoint}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (!response.ok) {
            if (response.status === 401) {
                localStorage.removeItem('token');
                showAuthForms();
            }
            const contentType = response.headers.get("content-type");
            if (contentType && contentType.indexOf("application/json") !== -1) {
                const errorData = await response.json();
                console.error(`HTTP error! status: ${response.status}, message: ${errorData.message}`);
                throw new Error(`HTTP error! status: ${response.status}, message: ${errorData.message}`);
            } else {
                const errorText = await response.text();
                console.error(`HTTP error! status: ${response.status}, message: ${errorText}`);
                throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
            }
        }
        const data = await response.json();
        console.log(`Data received from ${endpoint}:`, data);
        return data;
    }

    function logout() {
        localStorage.removeItem('token');
        showAuthForms();
    }
    
    document.getElementById('logout-button').addEventListener('click', logout);

    function renderPersons(data) {
        console.log('Rendering persons:', data);
        const list = document.getElementById('persons-list');
        list.innerHTML = '';
        data.forEach(person => {
            const li = document.createElement('li');
            li.textContent = `ID: ${person.id}, Name: ${person.name}, Age: ${person.age}, Gender: ${person.gender}, Address: ${person.address}`;
            list.appendChild(li);
        });
    }

    function renderPizzerias(data) {
        console.log('Rendering pizzerias:', data);
        const list = document.getElementById('pizzerias-list');
        list.innerHTML = '';
        data.forEach(pizzeria => {
            const li = document.createElement('li');
            li.textContent = `ID: ${pizzeria.id}, Name: ${pizzeria.name}`;
            list.appendChild(li);
        });
    }

    function renderVisits(data) {
        console.log('Rendering visits:', data);
        const list = document.getElementById('visits-list');
        list.innerHTML = ''; // Очищаем список перед добавлением новых элементов
    
        data.forEach(visit => {
            const li = document.createElement('li');
            li.textContent = `ID: ${visit.id}, Person Name: ${visit.name}, Pizzeria Name: ${visit.pizzeria_name}, Visit Date: ${visit.visit_date}`;
            list.appendChild(li);
        });
    }

    function renderOrders(data) {
        console.log('Rendering orders:', data);
        const list = document.getElementById('orders-list');
        list.innerHTML = '';
        data.forEach(order => {
            const li = document.createElement('li');
            li.textContent = `ID: ${order.id}, Person Name: ${order.name}, Menu ID: ${order.menu_id}, Order Date: ${order.order_date}`;
            list.appendChild(li);
        });
    }

    function renderMenus(data) {
        console.log('Rendering menus:', data);
        const list = document.getElementById('menus-list');
        list.innerHTML = '';
        data.forEach(menu => {
            const li = document.createElement('li');
            li.textContent = `ID: ${menu.id}, Pizzeria ID: ${menu.pizzeria_id}, Pizza Name: ${menu.pizza_name}, Price: ${menu.price}`;
            list.appendChild(li);
        });
    }

    document.getElementById('search-persons-form').addEventListener('submit', async (event) => {
        event.preventDefault();
        const city = document.getElementById('search-city').value;
        const minAge = document.getElementById('search-min-age').value;
    
        if (!city || !minAge) {
            alert('Please fill in all fields.');
            return;
        }
    
        try {
            const response = await fetch(`${baseURL}/api/persons/search?city=${city}&minAge=${minAge}`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
    
            if (response.ok) {
                const persons = await response.json();
                renderSearchResults(persons);
            } else {
                const errorData = await response.json();
                console.error('Failed to search persons:', errorData);
                alert('Failed to search persons. Please try again.');
            }
        } catch (error) {
            console.error('Error during persons search:', error);
            alert('Failed to search persons. Please try again.');
        }
    });
    
    function renderSearchResults(persons) {
        const list = document.getElementById('search-results');
        list.innerHTML = '';
        persons.forEach(person => {
            const li = document.createElement('li');
            li.textContent = `ID: ${person.id}, Name: ${person.name}, Age: ${person.age}, Gender: ${person.gender}, Address: ${person.address}`;
            list.appendChild(li);
        });
    }

    document.getElementById('average-age-form').addEventListener('submit', async (event) => {
        event.preventDefault();
        const city = document.getElementById('average-age-city').value;
    
        if (!city) {
            alert('Please fill in the city field.');
            return;
        }
    
        try {
            const response = await fetch(`${baseURL}/api/persons/average-age?city=${city}`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
    
            if (response.ok) {
                const data = await response.json();
                document.getElementById('average-age-result').textContent = `Average Age: ${data.average_age}`;
            } else {
                const errorData = await response.json();
                console.error('Failed to get average age:', errorData);
                alert('Failed to get average age. Please try again.');
            }
        } catch (error) {
            console.error('Error during average age request:', error);
            alert('Failed to get average age. Please try again.');
        }
    });
    
    document.getElementById('delete-older-than-form').addEventListener('submit', async (event) => {
        event.preventDefault();
        const maxAge = document.getElementById('delete-older-than-age').value;
    
        if (!maxAge) {
            alert('Please fill in the max age field.');
            return;
        }
    
        try {
            const response = await fetch(`${baseURL}/api/persons/delete-older-than?maxAge=${maxAge}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
    
            if (response.ok) {
                alert('Persons older than the specified age have been deleted.');
                loadData();
            } else {
                const errorData = await response.json();
                console.error('Failed to delete persons:', errorData);
                alert('Failed to delete persons. Please try again.');
            }
        } catch (error) {
            console.error('Error during delete persons request:', error);
            alert('Failed to delete persons. Please try again.');
        }
    });

    document.getElementById('update-person-with-trigger-form').addEventListener('submit', async (event) => {
        event.preventDefault();
        const id = document.getElementById('update-person-with-trigger-id').value;
        const name = document.getElementById('update-person-with-trigger-name').value;
        const age = document.getElementById('update-person-with-trigger-age').value;
        const gender = document.getElementById('update-person-with-trigger-gender').value;
        const address = document.getElementById('update-person-with-trigger-address').value;
    
        if (!id || !name || !age || !gender || !address) {
            alert('Please fill in all fields.');
            return;
        }
    
        try {
            const response = await fetch(`${baseURL}/api/persons?id=${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ name, age, gender, address })
            });
    
            if (response.ok) {
                alert('Person updated successfully with trigger!');
                loadData();
            } else {
                const errorData = await response.json();
                console.error('Failed to update person with trigger:', errorData);
                alert('Failed to update person with trigger. Please try again.');
            }
        } catch (error) {
            console.error('Error during person update with trigger:', error);
            alert('Failed to update person with trigger. Please try again.');
        }
    });
    
    async function loadData() {
        try {
            const persons = await fetchData('persons');
            renderPersons(persons);

            const pizzerias = await fetchData('pizzerias');
            renderPizzerias(pizzerias);

            const visits = await fetchData('visits');
            renderVisits(visits);

            const orders = await fetchData('orders');
            renderOrders(orders);

            const menus = await fetchData('menus');
            renderMenus(menus);
        } catch (error) {
            console.error('Error loading data:', error);
            alert('Failed to load data. Please try again.');
        }
    }

    document.getElementById('register').addEventListener('submit', async (event) => {
        event.preventDefault();
        const username = document.getElementById('register-username').value;
        const password = document.getElementById('register-password').value;

        if (!username || !password) {
            alert('Please fill in all fields.');
            return;
        }

        console.log('Registering user:', username);

        try {
            const response = await fetch(`${baseURL}/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });

            if (response.ok) {
                alert('Registration successful! Please login.');
            } else {
                const errorData = await response.json();
                console.error('Registration failed:', errorData);
                alert('Registration failed. Please try again.');
            }
        } catch (error) {
            console.error('Error during registration:', error);
            alert('Registration failed. Please try again.');
        }
    });

    document.getElementById('login').addEventListener('submit', async (event) => {
        event.preventDefault();
        const username = document.getElementById('login-username').value;
        const password = document.getElementById('login-password').value;

        if (!username || !password) {
            alert('Please fill in all fields.');
            return;
        }

        console.log('Logging in user:', username);

        try {
            const response = await fetch(`${baseURL}/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });

            if (response.ok) {
                const data = await response.json();
                localStorage.setItem('token', data.token);
                showContent();
                loadData();
            } else {
                const errorData = await response.json();
                console.error('Login failed:', errorData);
                alert('Login failed. Please try again.');
            }
        } catch (error) {
            console.error('Error during login:', error);
            alert('Login failed. Please try again.');
        }
    });

    // Добавление пользователя
    document.getElementById('create-person-form').addEventListener('submit', async (event) => {
        event.preventDefault();
        const name = document.getElementById('create-person-name').value;
        const age = document.getElementById('create-person-age').value;
        const gender = document.getElementById('create-person-gender').value;
        const address = document.getElementById('create-person-address').value;

        if (!name || !age || !gender || !address) {
            alert('Please fill in all fields.');
            return;
        }

        try {
            const response = await fetch(`${baseURL}/api/persons`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ name, age, gender, address })
            });

            if (response.ok) {
                alert('Person created successfully!');
                loadData();
            } else {
                const errorData = await response.json();
                console.error('Failed to create person:', errorData);
                alert('Failed to create person. Please try again.');
            }
        } catch (error) {
            console.error('Error during person creation:', error);
            alert('Failed to create person. Please try again.');
        }
    });

    // Обновление пользователя
    document.getElementById('update-person-form').addEventListener('submit', async (event) => {
        event.preventDefault();
        const id = document.getElementById('update-person-id').value;
        const name = document.getElementById('update-person-name').value;
        const age = document.getElementById('update-person-age').value;
        const gender = document.getElementById('update-person-gender').value;
        const address = document.getElementById('update-person-address').value;

        if (!id || !name || !age || !gender || !address) {
            alert('Please fill in all fields.');
            return;
        }

        try {
            const response = await fetch(`${baseURL}/api/persons?id=${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ name, age, gender, address })
            });

            if (response.ok) {
                alert('Person updated successfully!');
                loadData();
            } else {
                const errorData = await response.json();
                console.error('Failed to update person:', errorData);
                alert('Failed to update person. Please try again.');
            }
        } catch (error) {
            console.error('Error during person update:', error);
            alert('Failed to update person. Please try again.');
        }
    });

    // Удаление пользователя
    document.getElementById('delete-person-form').addEventListener('submit', async (event) => {
        event.preventDefault();
        const id = document.getElementById('delete-person-id').value;

        if (!id) {
            alert('Please fill in the ID field.');
            return;
        }

        try {
            const response = await fetch(`${baseURL}/api/persons?id=${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.ok) {
                alert('Person deleted successfully!');
                loadData();
            } else {
                const errorData = await response.json();
                console.error('Failed to delete person:', errorData);
                alert('Failed to delete person. Please try again.');
            }
        } catch (error) {
            console.error('Error during person deletion:', error);
            alert('Failed to delete person. Please try again.');
        }
    });
});
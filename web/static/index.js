$(document).ready(function () {
    let page = 1;
    const itemsPerPage = 10;
    const loadingIndicator = $('<div id="loading-indicator" class="text-center my-4"><div class="spinner-border" role="status"><span class="visually-hidden">Loading...</span></div></div>');

    function fetchData(page) {
        return $.ajax({
            url: '/buzz/content/list',
            method: 'POST',
            data: JSON.stringify({ PageIndex: page, PageSize: itemsPerPage }),
            headers: {
                'Authorization': 'Bearer ' + getToken(),
                'Content-Type': 'application/json'
            },
            success: function (data) {
                return data.contents;
            },
            error: function () {
                alert("Failed to fetch data");
            }
        });
    }

    function renderMasonry(contents) {
        const masonryContainer = $('#masonry');
        contents.forEach(item => {
            const cardHtml = `
                <div class="col-md-12 mb-4">
                    <div class="card d-flex flex-row" data-credit="${item.content_credit}">
                        <div class="card-img-left-wrapper">
                            <img src="${item.thumbnail}" class="card-img-left" alt="...">
                        </div>
                        <div class="card-body">
                            <h5 class="card-title">${item.title}</h5>
                            <p class="card-text">${item.platform}</p>
                        </div>
                    </div>
                </div>
            `;
            const cardElement = $(cardHtml);
            cardElement.click(function () {
                playAudio(item);
            });
            masonryContainer.append(cardElement);
        });
    }

    function loadMore() {
        $('body').append(loadingIndicator);
        fetchData(page).then(data => {
            renderMasonry(data.contents);
            page++;
            loadingIndicator.remove();
        });
    }

    function debounce(func, wait) {
        let timeout;
        return function (...args) {
            const later = () => {
                clearTimeout(timeout);
                func.apply(this, args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    $(window).scroll(debounce(function () {
        if ($(window).scrollTop() + $(window).height() >= $(document).height() - 100) {
            loadMore();
        }
    }, 200));

    // Initial load
    loadMore();

    var aplayer = new APlayer({
        container: document.getElementById('aplayer'),
        audio: {
            name: 'name',
            artist: 'artist',
            cover: 'https://fakeimg.pl/300x300',
            listFolded: true,
            theme: '#b7daff',
        }
    });

    function playAudio(data) {
        aplayer = new APlayer({
            container: document.getElementById('aplayer'),
            audio: {
                name: data.title,
                artist: data.text,
                url: `/buzz/content/stream/` + data.content_credit,
                cover: data.thumbnail,
                listFolded: true,
                theme: '#b7daff',
            }
        });
        aplayer.play();
    }

    // Back to top button
    const backToTopBtn = $('<button id="back-to-top" class="btn btn-primary">Top</button>');
    $('body').append(backToTopBtn);
    backToTopBtn.hide();

    $(window).scroll(function () {
        if ($(window).scrollTop() > 300) {
            backToTopBtn.fadeIn();
        } else {
            backToTopBtn.fadeOut();
        }
    });

    backToTopBtn.click(function () {
        $('html, body').animate({ scrollTop: 0 }, '300');
    });

    // Function to save JWT token
    function saveToken(token) {
        localStorage.setItem('jwtToken', token);
    }

    // Function to get JWT token
    function getToken() {
        return localStorage.getItem('jwtToken');
    }

    // Function to remove JWT token
    function removeToken() {
        localStorage.removeItem('jwtToken');
    }

    // Function to update login state
    function updateLoginState() {
        const token = getToken();
        if (token) {
            $.ajax({
                url: '/auth/current_user',
                method: 'GET',
                headers: {
                    'Authorization': 'Bearer ' + token
                },
                success: function (data) {
                    $('#login-btn').hide();
                    $('#user-info').removeClass('d-none');
                    $('#username').text(data.username);
                },
                error: function () {
                    removeToken();
                    $('#login-btn').show();
                    $('#user-info').addClass('d-none');
                }
            });
        } else {
            $('#login-btn').show();
            $('#user-info').addClass('d-none');
        }
    }

    // Login button click event
    $('#login-btn').click(function () {
        $('#loginModal').modal('show');
    });

    // Submit login button click event
    $('#submit-login').click(function () {
        const username = $('#login-username').val();
        const password = $('#login-password').val();
        $.ajax({
            url: '/auth/login',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ username: username, password: password }),
            success: function (data) {
                saveToken(data.token);
                updateLoginState();
                $('#loginModal').modal('hide');
            },
            error: function () {
                alert("Invalid credentials");
            }
        });
    });

    // Logout button click event
    $('#logout-btn').click(function () {
        const token = getToken();
        if (token) {
            $.ajax({
                url: '/auth/logout',
                method: 'POST',
                headers: {
                    'Authorization': 'Bearer ' + token
                },
                success: function () {
                    removeToken();
                    updateLoginState();
                },
                error: function () {
                    alert("Failed to logout");
                }
            });
        }
    });

    // Initial login state update
    updateLoginState();

    $('#subscribe-btn').click(function () {
        $('#subscribeModal').modal('show');
    });

    $('#submit-subscription').click(function () {
        const url = $('#subscription-url').val();
        if (url) {
            console.log("Subscribed to:", url);
            $('#subscribeModal').modal('hide');
        } else {
            alert("Please enter a URL");
        }
    });

    // Initial login state update
    updateLoginState();
});

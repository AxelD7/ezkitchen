var acc = document.getElementsByClassName("accordion");
for (var i = 0; i < acc.length; i++) {
    acc[i].addEventListener("click", function () {
        for (var j = 0; j < acc.length; j++) {
                acc[j].classList.remove("active");
                acc[j].nextElementSibling.style.display = "none";
        }
        this.classList.toggle("active");
        this.nextElementSibling.style.display = "block";
    });
}

const button = document.querySelector('#menu-button');
const menu1 = document.querySelector('#menu-links-1');
const menu2 = document.querySelector('#menu-links-2');
const navbar = document.querySelector('#navbar');
document.getElementById("animated-bg-video").removeAttribute("autoplay");
let memoizedPosition = 'top';

button.addEventListener('click', () => {
	menu1.classList.toggle('hidden');
	menu2.classList.toggle('hidden');
});

const getPositionFromScroll = position => {
	return position > 10 ? 'bottom' : 'top'
}

const toggleNavbar = () => {
	navbar.classList.toggle('bg-transparent')
	navbar.classList.toggle('bg-white')
	navbar.classList.toggle('text-white')
}

document.querySelector('#wrapper').addEventListener('scroll', ({ target: { scrollTop } }) => {
	if (getPositionFromScroll(scrollTop) !== memoizedPosition) {
		memoizedPosition = getPositionFromScroll(scrollTop);
		toggleNavbar();
	}
});

const registerVideo = (bound, video) => {
	notEnabled = true;
	cds = document.getElementsByClassName('cd')
	bound = document.querySelector(bound);
	video = document.querySelector(video);
	const scrollVideo = ()=>{
		if(video.duration) {
			const distanceFromTop = window.scrollY + bound.getBoundingClientRect().top;
			const rawPercentScrolled = (window.scrollY - distanceFromTop) / (bound.scrollHeight - window.innerHeight);
			const percentScrolled = Math.min(Math.max(rawPercentScrolled, 0), 1);
			
			if (notEnabled && percentScrolled > 0.1) {
				notEnabled = false;
				Array.from(cds).forEach((element, index) => {
					setTimeout(() => {
						element.classList.remove('opacity-0');
					}, 100*index);
				});
			}
			video.currentTime = video.duration * percentScrolled;
		}
		requestAnimationFrame(scrollVideo);
	}
	requestAnimationFrame(scrollVideo);
}

registerVideo("#bound-one", "#bound-one video");
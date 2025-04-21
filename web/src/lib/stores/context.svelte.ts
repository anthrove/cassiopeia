import { writable } from 'svelte/store';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

const context = writable('');

if (browser) {
	const url = new URL(window.location.href);
	const initial = url.searchParams.get('tenant') ?? '';
	context.set(initial);

    // push to url
    context.subscribe((value)=>{
        const url = new URL(window.location.href);
        if (url.searchParams.get('tenant') !== value) {
            if(value){
                url.searchParams.set('tenant', value);
            }else{
                url.searchParams.delete('tenant')
            }
            goto(`${url.pathname}?${url.searchParams.toString()}`, {
                replaceState: true,
                noScroll: true
            });
        }
    })

	// Handle back/forward button
	window.addEventListener('popstate', () => {
		const url = new URL(window.location.href);
		context.set(url.searchParams.get('tenant') ?? '');
	});
}

export { context };

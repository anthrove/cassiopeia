import { writable, type Writable } from 'svelte/store';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

const context = writable('');

if (browser) {
    const url = new URL(window.location.href);
    const initial = url.searchParams.get('tenant') ?? '';
    context.set(initial);
}

let contextReference = $state<{[key:string]:Writable<string>}>({})
function create(parameter:string){
    // this is a singleton construct
    if(contextReference[parameter]){
        return contextReference[parameter]
    }

    // create
    const url = new URL(window.location.href);
    const initial = url.searchParams.get(parameter) ?? '';
    contextReference[parameter] = writable(initial)

    // push changes to url
    contextReference[parameter].subscribe((value)=>{
        const url = new URL(window.location.href);
        if (url.searchParams.get(parameter) !== value) {
            if(value){
                url.searchParams.set(parameter, value);
            }else{
                url.searchParams.delete(parameter)
            }
            goto(`${url.pathname}?${url.searchParams.toString()}`, {
                //replaceState: true,
                noScroll: true
            });
        }
    })

    return contextReference[parameter]
}


// Handle back/forward button
window.addEventListener('popstate', () => {
    const url = new URL(window.location.href);
    for(const key in Object.keys(contextReference)){
        contextReference[key].set(url.searchParams.get(key) ?? '');
    }
});

export { create };

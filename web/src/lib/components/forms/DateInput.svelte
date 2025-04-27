<script lang="ts">
    import { DatePicker } from "@svelte-plugins/datepicker";

    let startDate = $state(new Date());
    let dateFormat = "MM/dd/yy";
    let isOpen = $state(false);

    const toggleDatePicker = () => (isOpen = !isOpen);

    let {
        value = $bindable('')
    } = $props()

    function toDatetimeLocalString(date) {
        if (!date) return "";

        const pad = (n) => n.toString().padStart(2, "0");

        const year = date.getFullYear();
        const month = pad(date.getMonth() + 1); // Months are 0-indexed
        const day = pad(date.getDate());
        const hours = pad(date.getHours());
        const minutes = pad(date.getMinutes());

        return `${year}-${month}-${day}T${hours}:${minutes}`;
    }

    const formatDate = (dateString: string) => {
        if (isNaN(new Date(dateString))) {
            return "";
        }

        return toDatetimeLocalString(new Date(dateString));
    };
    let formattedStartDate = $state(formatDate(startDate));

    const onChange = () => {
        startDate = new Date(formattedStartDate);
    };

    $effect(()=>{
        formattedStartDate = formatDate(startDate);
        value = new Date(startDate).toISOString()
    })
    /*
    $effect(()=>{
        if(formattedStartDate != value){
            formattedStartDate = value
        }
    })*/
    //$: formattedStartDate = formatDate(startDate);
</script>

<DatePicker bind:isOpen bind:startDate showTimePicker enableFutureDates enablePastDates={false}>
    <input
        type="text"
        placeholder="Select date"
        bind:value={formattedStartDate}
        onclick={toggleDatePicker}
    />
</DatePicker>

<style>

</style>

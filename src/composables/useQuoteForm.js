import { reactive, computed } from "vue";

export function useQuoteForm() {
    const form = reactive({
        appliance: "",
        name: "",
        phone: "",
        zipCode: "",
        description: "",
    });

    const isFormValid = computed(() =>
        form.appliance && form.name && form.phone
    );

    function resetForm() {
        Object.keys(form).forEach(k => (form[k] = ""));
    }

    return {
        form,
        isFormValid,
        resetForm,
    };
}

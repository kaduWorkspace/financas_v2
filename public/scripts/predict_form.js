import { DateUtils, CurrencyUtils, FormUtils, StorageUtils } from './inputUtils.js';
// Check if the script has already been loaded
if (typeof formValidationInitialized === 'undefined') {
    // Function to format and validate number inputs (handles both comma and dot decimals)
    function parseNumber(inputValue) {
        const normalizedValue = inputValue.replace(/\./g, '').replace(',', '.');
        const number = parseFloat(normalizedValue);
        return isNaN(number) ? null : number;
    }

    // Function to validate individual fields
    function validateField(fieldId, isRequired = true) {
        const input = document.getElementById(fieldId);
        const errorElement = document.getElementById(`error_${fieldId}`);
        const value = parseNumber(input.value);

        if (isRequired && (!value && value !== 0)) {
            errorElement.textContent = 'Este campo é obrigatório';
            errorElement.classList.remove('hidden');
            return false;
        }

        if (value !== null && value < 0) {
            errorElement.textContent = 'O valor não pode ser negativo';
            errorElement.classList.remove('hidden');
            return false;
        }

        errorElement.classList.add('hidden');
        return true;
    }

    // Main validation function
    function validateForm() {
        const isValorFuturoValid = validateField('valor_futuro');
        const isTaxaJurosValid = validateField('taxa_juros_anual');
        const isValorInicialValid = validateField('valor_inicial', false);

        return isValorFuturoValid && isTaxaJurosValid && isValorInicialValid;
    }

    // Set up event listeners for real-time validation
    function setupFormValidation() {
        const form = document.getElementById('formulario_prever');
        if (!form) return;
        form.addEventListener("input", e => {
            const {target} = e
            if(["taxa_juros_anual", "valor_futuro", "valor_inicial"].includes(target.id)) return CurrencyUtils.handleCurrencyInput(e)
        })
        // Validate on input change
        document.getElementById('valor_futuro').addEventListener('input', () => validateField('valor_futuro'));
        document.getElementById('taxa_juros_anual').addEventListener('input', () => validateField('taxa_juros_anual'));
        document.getElementById('valor_inicial').addEventListener('input', () => validateField('valor_inicial', false));

        // Validate before HTMX request
        form.addEventListener('htmx:beforeRequest', function(event) {
            if (!validateForm()) {
                event.preventDefault();
                // Scroll to first error
                const firstError = document.querySelector('[id^="error_"]:not(.hidden)');
                if (firstError) {
                    firstError.scrollIntoView({ behavior: 'smooth', block: 'center' });
                }
            } else {
                // Format values for submission
                document.getElementById('valor_futuro_input').value = parseNumber(document.getElementById('valor_futuro').value);
                document.getElementById('taxa_juros_anual_input').value = parseNumber(document.getElementById('taxa_juros_anual').value);
                document.getElementById('valor_inicial_input').value = parseNumber(document.getElementById('valor_inicial').value) || 0;
            }
        });
    }

    setupFormValidation();

    var formValidationInitialized = true;
}

document.addEventListener('DOMContentLoaded', () => {
    const generateBtn = document.getElementById('generate-btn');
    const generateResult = document.getElementById('generate-result');

    const addressBtn = document.getElementById('address-btn');
    const addressResult = document.getElementById('address-result');

    const signBtn = document.getElementById('sign-btn');
    const signMessage = document.getElementById('sign-message');
    const signResult = document.getElementById('sign-result');

    const verifyBtn = document.getElementById('verify-btn');
    const verifyMessage = document.getElementById('verify-message');
    const verifySignature = document.getElementById('verify-signature');
    const verifyResult = document.getElementById('verify-result');

    const transactionBtn = document.getElementById('transaction-btn');
    const txToAddress = document.getElementById('tx-to-address');
    const txValue = document.getElementById('tx-value');
    const transactionResult = document.getElementById('transaction-result');

    generateBtn.addEventListener('click', async () => {
        const response = await fetch('/generate');
        const data = await response.json();
        generateResult.textContent = `Private Key: ${data.private_key}, Address: ${data.address}`;
    });

    addressBtn.addEventListener('click', async () => {
        const response = await fetch('/address');
        const data = await response.json();
        addressResult.textContent = `Address: ${data.address}`;
    });

    signBtn.addEventListener('click', async () => {
        const message = signMessage.value;
        const response = await fetch('/sign', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ message })
        });
        const data = await response.json();
        signResult.textContent = `Signature: ${data.signature}`;
    });

    verifyBtn.addEventListener('click', async () => {
        const message = verifyMessage.value;
        const signature = verifySignature.value;
        const response = await fetch('/verify', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ message, signature })
        });
        const data = await response.json();
        verifyResult.textContent = `Valid: ${data.valid}`;
    });

    transactionBtn.addEventListener('click', async () => {
        const toAddress = txToAddress.value;
        const value = txValue.value;
        const response = await fetch('/transaction', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ to_address: toAddress, value: parseInt(value) })
        });
        const data = await response.json();
        transactionResult.textContent = `Transaction Hash: ${data.transaction_hash}`;
    });
});

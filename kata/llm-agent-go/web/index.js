// index.js

const { createApp } = Vue;

document.addEventListener('DOMContentLoaded', () => {
    createApp({
        data() {
            return {
                answer: '',
                prompts: [],
                formData: {
                    prompt: {
                        name: '',
                        system_prompt: '',
                        user_prompt: '',
                        output_language: 'chinese'
                    },
                    code_path: './cmd/web.go',
                    stream: true
                }
            };
        },
        methods: {
            async fetchCommands() {
                try {
                    const response = await fetch('/api/v1/commands');
                    if (!response.ok) {
                        throw new Error(`HTTP error! Status: ${response.status}`);
                    }
                    const result = await response.json();
                    this.prompts = result;
                    if (this.prompts.length > 0) {
                        this.formData.prompt.name = this.prompts[0].name;
                        this.formData.prompt.system_prompt = this.prompts[0].system_prompt.replace(/\\n/g, '\n');
                        this.formData.prompt.user_prompt = this.prompts[0].user_prompt.replace(/\\n/g, '\n');
                    }
                } catch (error) {
                    alert('Failed to load command list: ' + error.message);
                }
            },
            async submitRequest() {
                this.answer = 'Loading...';

                if (this.formData.stream) {
                    this.startStreaming();
                    return;
                }

                try {
                    const response = await fetch('/api/v1/process', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify(this.formData)
                    });

                    if (!response.ok) {
                        throw new Error(`HTTP error! Status: ${response.status}`);
                    }

                    const result = await response.json();
                    this.answer = result.answer;
                } catch (error) {
                    this.answer = 'Error: ' + error.message;
                }
            },
            async startStreaming() {
                const ws = new WebSocket(`ws://${window.location.host}/api/v1/stream`);

                const formData = JSON.parse(JSON.stringify(this.formData)); // Deep copy

                ws.onopen = () => {
                    console.log('WebSocket opened, sending data:', formData);
                    ws.send(JSON.stringify(formData));
                };

                ws.onmessage = (event) => {
                    const chunk = event.data;
                    this.answer += chunk;
                };

                ws.onerror = (err) => {
                    console.error('WebSocket Error:', err);
                    this.answer += '\n[Error occurred]';
                    ws.close();
                };

                ws.onclose = () => {
                    console.log('WebSocket closed');
                };
            }
        },
        watch: {
            'formData.prompt.name': function (newName) {
                const selected = this.prompts.find(p => p.name === newName);
                if (selected) {
                    this.formData.prompt.system_prompt = selected.system_prompt.replace(/\\n/g, '\n');
                    this.formData.prompt.user_prompt = selected.user_prompt.replace(/\\n/g, '\n');
                }
            }
        },
        mounted() {
            this.fetchCommands(); // Load commands when component is mounted
        }
    }).mount('#app');
});
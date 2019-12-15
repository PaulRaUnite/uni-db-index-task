<template>
    <div class="root">
        <va-button :type="type" :size="size" :active="active"
                   :disabled="disabled" :block="block" :loading="loading"
                   :round="round" :focused="focused" :tall="tall"
                   @click="onClick">
            <transition appear>
                <va-icon v-if="state === 'success'" type="check" key="success"/>
                <va-icon v-if="state === 'failure'" type="times" key="failure"/>
                <div v-if="state === 'default' || state=== 'loading'" class="root" key="default">
                    <slot/>
                </div>
            </transition>
        </va-button>
    </div>
</template>

<script>
    export default {
        name: "ctx-button",
        data() {
            return {
                loading: false,
            }
        },
        props: {
            type: {
                type: String,
                default: 'default',
                required: false,
                validator(v) {
                    return [
                        'default',
                        'primary',
                        'primary-light',
                        'primary-dark',
                        'paleblue',
                        'success',
                        'info',
                        'warning',
                        'danger',
                        'subtle',
                        'link',
                        'subtle-link',
                        'active',
                        'dark',
                        'darker',
                        'help',
                        'help-light',
                        'help-dark',
                        'black'
                    ].includes(v)
                }
            },
            size: {
                type: String,
                default: 'md',
                required: false,
                validator(v) {
                    return ['xs', 'sm', 'md', 'lg'].includes(v)
                }
            },
            active: {
                type: Boolean,
                default: false,
                required: false
            },
            disabled: {
                type: Boolean,
                default: false,
                required: false
            },
            block: {
                type: Boolean,
                default: false,
                required: false
            },
            round: {
                type: Boolean,
                default: false,
                required: false
            },
            focused: {
                type: Boolean,
                default: false,
                required: false
            },
            tall: {
                type: Boolean,
                default: false,
                required: false
            },
            iconBefore: {
                type: String,
                required: false
            },
            iconAfter: {
                type: String,
                required: false
            },
            state: {
                type: String,
                default: 'default',
                required: false,
                validator(v) {
                    return [
                        'default',
                        'success',
                        'failure',
                        'loading'
                    ].includes(v)
                }
            },
        },
        watch: {
            state() {
                this.loading = this.state === "loading";
            }
        },
        methods: {
            onClick(event) {
                if (this.disabled) {
                    return
                }
                this.$emit('click', event)
            }
        }
    }
</script>

<style scoped>
    .root {
        display: inline-block;
    }
</style>
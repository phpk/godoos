import { reactive, ref, watch } from 'vue'

export type ScreenRecorderStatus =
    | 'recording'
    | 'idle'
    | 'error'
    | 'stopped'
    | 'paused'
    | 'permission-requested'

export interface useScreenRecorderParams {
    options?: MediaRecorderOptions
    audio?: boolean
}

const useScreenRecorder = (params: useScreenRecorderParams = {}) => {
    const { options, audio = false } = params
    const blobUrlRef = ref<string>()
    const blobRef = ref<Blob>()
    const errorRef = ref()
    const mediaRecorderRef = ref<MediaRecorder>()
    const statusRef = ref<ScreenRecorderStatus>('permission-requested')
    const streamsRef = ref<{
        audio?: MediaStreamTrack
        screen?: MediaStreamTrack
    }>({ audio: undefined, screen: undefined })

    watch(mediaRecorderRef, () => {
        if (mediaRecorderRef.value) {
            mediaRecorderRef.value.ondataavailable = (event) => {
                const url = window.URL.createObjectURL(event.data)
                blobUrlRef.value = url
                blobRef.value = event.data
            }
        }
    })

    const requestMediaStream = async () => {
        try {
            // @ts-ignore
            const displayMedia = await navigator.mediaDevices.getDisplayMedia()

            let userMedia: MediaStream

            if (audio) {
                userMedia = await navigator.mediaDevices.getUserMedia({ audio })
            }

            // @ts-ignore
            const tracks = [...displayMedia?.getTracks(), ...(userMedia?.getTracks() || [])]
            if (tracks) statusRef.value = 'idle'
            const stream: MediaStream = new MediaStream(tracks)
            const mediaRecorder = new MediaRecorder(stream, options)
            mediaRecorderRef.value = mediaRecorder

            streamsRef.value = {
                audio:
                    // @ts-ignore
                    userMedia?.getTracks().find(track => track.kind === 'audio'),
                screen:
                    displayMedia
                        .getTracks()
                        .find((track: MediaStreamTrack) => track.kind === 'video')
            }

            return mediaRecorder
        } catch (e) {
            errorRef.value = e
            statusRef.value = 'error'
        }
    }

    const stopRecording = () => {
        if (!mediaRecorderRef.value) throw new Error('No media stream!')
        mediaRecorderRef.value.stop()

        statusRef.value = 'stopped'

        mediaRecorderRef.value.stream.getTracks().forEach((track) => {
            track.stop()
        })

        mediaRecorderRef.value = undefined
    }

    const startRecording = async () => {
        if (!mediaRecorderRef.value) {
            mediaRecorderRef.value = await requestMediaStream()
        }

        if (mediaRecorderRef.value) {
            mediaRecorderRef.value?.start()
            statusRef.value = 'recording'
        }
    }

    const pauseRecording = () => {
        if (!mediaRecorderRef.value) throw new Error('No media stream!')
        mediaRecorderRef.value?.pause()
        statusRef.value = 'paused'
    }

    const resumeRecording = () => {
        if (!mediaRecorderRef.value) throw new Error('No media stream!')
        mediaRecorderRef.value?.resume()
        statusRef.value = 'recording'
    }

    const resetRecording = () => {
        blobUrlRef.value = undefined
        blobRef.value = undefined
        errorRef.value = undefined
        mediaRecorderRef.value = undefined
        statusRef.value = 'idle'
    }

    return reactive({
        blob: blobRef,
        blobUrl: blobUrlRef,
        error: errorRef,
        status: statusRef,
        streams: streamsRef,
        actions: {
            startRecording,
            stopRecording,
            pauseRecording,
            resetRecording,
            resumeRecording
        }
    })
}

export default useScreenRecorder
@Library('dst-shared@master') _

dockerBuildPipeline {
        githubPushRepo = "Cray-HPE/hms-bmc-networkprotocol"
        repository = "cray"
        imagePrefix = "hms"
        app = "networkprotocol"
        name = "hms-bmc-networkprotocol"
        description = "Cray HMS common BMC network protocol control."
        dockerfile = "Dockerfile"
        slackNotification = ["", "", false, false, true, true]
        product = "internal"
}


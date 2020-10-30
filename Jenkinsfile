@Library('dst-shared@master') _

dockerBuildPipeline {
        repository = "cray"
        imagePrefix = "hms"
        app = "networkprotocol"
        name = "hms-bmc-networkprotocol"
        description = "Cray HMS common BMC network protocol control."
        dockerfile = "Dockerfile"
        slackNotification = ["", "", false, false, true, true]
        product = "internal"
}


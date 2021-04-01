const config = require("./config.js");
const fs = require("fs");
const Discord = require("discord.js");
const low = require("lowdb");
const FileSync = require("lowdb/adapters/FileSync");
const serialize = require("serialize-javascript");
const client = new Discord.Client({
    partials: ["MESSAGE", "USER", "REACTION", "GUILD_MEMBER"],
});
const admins = new Set([
    "410635076842422293", // Iman
    "191587255557554177", // Oppy
]);
const advertiserRoles = ["Heroic Advertiser"];
const eliteAdvertiserWeights = {
    "Advertiser Trainer": 1,
    "Support Team": 1,
    "Elite Advertiser": 1,
    "Supreme Advertiser": 1.75,
    "Monster Advertiser": 2.5,
    "Titan Advertiser": 3.25,
    "Legendary Advertiser": 4,
    "Demigod Advertiser": 5,
    "Pantheon Advertiser": 5,
    "The Eternal Advertiser": 5,
};
const reactionArray = ["👍", "⏭"];
const boostRequestsBySignupMessageId = new Map();
const boostRequestTimeouts = new Map();

const db = low(
    new FileSync("db.json", {
        defaultValue: {
            instantBoostRequestCredits: {},
        },
    })
);

client.login(config.token);
let areBoostRequestsLoaded = false;
client.on("ready", () => {
    console.log("client online!");
    if (!areBoostRequestsLoaded) {
        areBoostRequestsLoaded = true;
        loadBoostRequests();
    }
});

function loadBoostRequests() {
    fs.readFile(
        `${__dirname}/boost-requests.js`,
        (err, loadedBoostRequests) => {
            if (err) {
                console.error("Failed to load boost requests", err);
                return;
            }
            try {
                const boostRequests = eval(
                    `(${loadedBoostRequests.toString()})`
                );
                boostRequests.forEach((boostRequest, key) => {
                    // XXX Backwards compatibility, remove after all requests expire (3 days)
                    if (!boostRequest.queuedAdvertisers) {
                        boostRequest.queuedAdvertisers = [];
                    }
                    boostRequestsBySignupMessageId.set(key, boostRequest);
                    addTimers(boostRequest);
                });
            } catch (err) {
                console.error("Failed to parse boost requests", err);
            }
        }
    );
}

function saveBoostRequests() {
    const serialziedBoostRequests = serialize(boostRequestsBySignupMessageId, {
        unsafe: true,
        ignoreFunction: true,
    });
    fs.writeFileSync(`${__dirname}/boost-requests.js`, serialziedBoostRequests);
}

// TODO: Edit first embed to get of 2nd embed
// TODO: Clean up embeds
// TODO: Fix bug so it doesnt respond to one person and mark all with checks
// TODO: Put every reaction.message.id into an array and check again the current one being operated on so we know if we should be using it or not

// Event Catcher when users react to messages
client.on("messageReactionAdd", async (reaction, user) => {
    try {
        if (reaction.partial) {
            await reaction.fetch();
        }
        if (user.partial) {
            await user.fetch();
        }
    } catch (err) {
        console.error(reaction.id, user.id, err);
        return;
    }
    const boostRequest = boostRequestsBySignupMessageId.get(
        reaction.message.id
    );
    const guildMember = await reaction.message.guild.members.fetch(user);

    if (boostRequest && !user.bot) {
        const isAdvertiser = guildMember.roles.cache.some((role) =>
            advertiserRoles.includes(role.name)
        );
        const eliteAvertiserRole = guildMember.roles.cache.reduce(
            (best, current) => {
                if (current.name in eliteAdvertiserWeights) {
                    return (eliteAdvertiserWeights[best] ?? -Infinity) >
                        eliteAdvertiserWeights[current.name]
                        ? best
                        : current.name;
                }
                return best;
            },
            null
        );
        if (
            reaction.emoji.name === "👍" ||
            (reaction.emoji.name == "⏭" &&
                ((isAdvertiser && boostRequest.isClaimableByAdvertisers) ||
                    (eliteAvertiserRole &&
                        boostRequest.isClaimableByEliteAdvertisers)))
        ) {
            if (
                (boostRequest.isClaimableByAdvertisers && isAdvertiser) ||
                (boostRequest.isClaimableByEliteAdvertisers &&
                    eliteAvertiserRole)
            ) {
                await setWinner(reaction.message, user);
            } else if (isAdvertiser || eliteAdvertiserWeights) {
                boostRequest.queuedAdvertisers.push({
                    id: user.id,
                    isElite: eliteAvertiserRole !== null,
                    weight: eliteAdvertiserWeights[eliteAvertiserRole],
                });
            }
        } else if (reaction.emoji.name == "⏭") {
            let availableCredits = db
                .get("instantBoostRequestCredits")
                .get([user.id], 0);
            if (availableCredits > 0) {
                availableCredits--;
                db.get("instantBoostRequestCredits")
                    .set([user.id], availableCredits)
                    .write();
                const dmChannel = user.dmChannel ?? (await user.createDM());
                await setWinner(reaction.message, user);
                try {
                    await dmChannel.send(
                        `You used an instant boost request. You have ${availableCredits} credits remaining.`
                    );
                } catch (err) {
                    // they're blocking dms
                }
            }
        }
    }
});

client.on("messageReactionRemove", async (reaction, user) => {
    try {
        if (reaction.partial) {
            await reaction.fetch();
        }
    } catch (err) {
        return;
    }
    const boostRequest = boostRequestsBySignupMessageId.get(
        reaction.message.id
    );
    if (boostRequest) {
        boostRequest.queuedAdvertisers = boostRequest.queuedAdvertisers.filter(
            (advertiser) => advertiser.id !== user.id
        );
    }
});

// Event Catcher when users send a message
client.on("message", async (message) => {
    try {
        if (message.partial) {
            await message.fetch();
        }
    } catch (err) {
        return;
    }
    if (message.author.equals(client.user)) return;
    console.log(message.content);
    const boostRequestChannel = config.boostRequestChannelId.find(
        (chan) => chan.id == message.channel.id
    );
    // If User is not a bot AND is messsaging in BoostRequest Channel
    if (
        boostRequestChannel &&
        (!message.author.bot || !boostRequestChannel.notifyBuyer)
    ) {
        // Create embed in the Backend Channel
        if (!boostRequestChannel.useBuyerMessage) {
            if (!(await sendBuyerWaitingMessage(message))) {
                return;
            }
        }
        const signupMessage = boostRequestChannel.useBuyerMessage
            ? message
            : await BREmbed(message, boostRequestChannel.backendId);
        const reactPromises = reactionArray.map((emoji) =>
            signupMessage.react(emoji)
        );
        await Promise.all(reactPromises);

        const buyerDiscordName =
            message.embeds.length >= 1
                ? message.embeds[0].fields.find((field) =>
                      field.name.toLowerCase().includes("battletag")
                  )?.value
                : undefined;
        const boostRequest = {
            channelId: message.channel.id,
            requesterId: message.author.id,
            createdAt: message.createdAt,
            backendChannelId: boostRequestChannel.backendId,
            buyerDiscordName: buyerDiscordName,
            isClaimableByAdvertisers: false,
            queuedAdvertisers: [],
            signupMessageId: signupMessage.id,
            message: message.content,
        };
        addTimers(boostRequest);
        boostRequestsBySignupMessageId.set(signupMessage.id, boostRequest);
    } else {
        // Command
        const [command, ...args] = message.content.split(" ");
        if (
            command == "!boostrequestcredit" &&
            args.length >= 2 &&
            admins.has(message.author.id)
        ) {
            const [userQuery, credits] = args;
            const user = await findUser(userQuery, message.guild.id);
            if (!user) {
                message.reply(`User "${userQuery}" not found`);
                return;
            }
            let numCredits;
            try {
                numCredits = parseInt(credits);
            } catch (err) {
                message.reply(`Invalid integer: ${numCredits}`);
            }
            const existingCredits = db
                .get("instantBoostRequestCredits")
                .get([user.id], 0)
                .value();
            const newCredits = existingCredits + numCredits;
            db.get("instantBoostRequestCredits")
                .set([user.id], newCredits)
                .write();
            await message.reply(
                `${user.tag} now has ${newCredits} instant boost request credits.`
            );
            const dmChannel = user.dmChannel ?? (await user.createDM());
            try {
                await dmChannel.send(
                    `You were granted ${numCredits} instant boost request credits. You have ${newCredits} total credits.`
                );
            } catch (err) {
                // they're blocking dms
            }
        }
    }
});

async function sendBuyerWaitingMessage(message) {
    const embed = new Discord.MessageEmbed()
        .setTitle("Huokan Boosting Community Boost Request")
        .setDescription(message.content)
        .setThumbnail(message.author.avatarURL())
        .setAuthor(`${message.author.username}#${message.author.discriminator}`)
        .setFooter(
            "Huokan Boosting Community",
            "https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png"
        )
        .setTimestamp();
    try {
        const dmChannel =
            message.author.dmChannel ?? (await message.author.createDM());
        await dmChannel.send(
            "Please wait while we find an advertiser to complete your request.",
            embed
        );
        if (message.deletable) {
            await message.delete();
        }
    } catch (err) {
        if (err.code === 50007) {
            // Cannot send messages to this user
            const reply = await message.reply(
                'I can\'t DM you! Please allow DMs from server members by right clicking the server and enabling "Allow direct messages from server members." in Privacy Settings, and then post your message again.'
            );
            setTimeout(() => {
                message.delete().catch(() => {
                    // ignore
                });
                reply.delete().catch(() => {
                    // ignore
                });
            }, 30000);
        } else {
            console.error(err);
        }
        return false;
    }
    return true;
}

function addTimers(boostRequest) {
    boostRequestTimeouts.set(boostRequest, [
        setTimeout(async () => {
            const advertisers = boostRequest.queuedAdvertisers.filter(
                (advertiser) => advertiser.isElite
            );
            if (advertisers.length >= 1) {
                const winner = getRandomAdvertiserWeighted(advertisers);
                try {
                    const user = await client.users.fetch(winner.id);
                    const channel = await client.channels.fetch(
                        boostRequest.backendChannelId
                    );
                    const signupMessage = await channel.messages.fetch(
                        boostRequest.signupMessageId
                    );
                    await setWinner(signupMessage, user);
                } catch (err) {
                    console.error(err);
                    boostRequest.isClaimableByEliteAdvertisers = true;
                }
            } else {
                boostRequest.isClaimableByEliteAdvertisers = true;
            }
        }, Math.max(0, 20000 - (new Date() - boostRequest.createdAt))),
        setTimeout(async () => {
            if (boostRequest.queuedAdvertisers.length >= 1) {
                try {
                    const chosenAdvertiser =
                        boostRequest.queuedAdvertisers[
                            Math.floor(
                                Math.random() *
                                    boostRequest.queuedAdvertisers.length
                            )
                        ];
                    const userId = chosenAdvertiser.id;
                    const user = await client.users.fetch(userId);
                    const channel = await client.channels.fetch(
                        boostRequest.backendChannelId
                    );
                    const signupMessage = await channel.messages.fetch(
                        boostRequest.signupMessageId
                    );
                    await setWinner(signupMessage, user);
                } catch (err) {
                    console.error(err);
                    boostRequest.isClaimableByAdvertisers = true;
                }
            } else {
                boostRequest.isClaimableByAdvertisers = true;
            }
        }, Math.max(0, 60000 - (new Date() - boostRequest.createdAt))),
        // 1 minute
        setTimeout(() => {
            console.log("Deleting expired boost request.");
            boostRequestsBySignupMessageId.delete(boostRequest.signupMessageId);
        }, 259200000),
        // 72 hours
    ]);
}

async function setWinner(message, winner) {
    const signupMessage = boostRequestsBySignupMessageId.get(message.id);
    if (!signupMessage) {
        return;
    }
    boostRequestsBySignupMessageId.delete(message.id);
    boostRequestTimeouts
        .get(signupMessage)
        .forEach((timeoutId) => clearTimeout(timeoutId));
    const boostRequestChannel = config.boostRequestChannelId.find(
        (chan) => chan.id == signupMessage.channelId
    );

    try {
        const winnerName = winner.username;
        console.log(winnerName + " won!");
        // remove reactions.
        try {
            await message.reactions.removeAll();
        } catch (err) {
            console.error("Failed to clear reactions: ", err);
        }
        await message.react("✅");
        await sendEmbed(winner, signupMessage, boostRequestChannel);
    } catch (error) {
        console.error("One of the emojis failed to react.", error);
    }
}

async function BREmbed(brMessage, channelId) {
    // Variable to eaily add hyperlink to the original message.
    const exampleEmbed = new Discord.MessageEmbed()
        .setColor("#0000FF")
        .setTitle("New Boost Request")
        .setTimestamp()
        .setFooter(
            "Huokan Boosting Community",
            "https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png"
        );
    if (brMessage.embeds.length >= 1) {
        exampleEmbed.addFields(brMessage.embeds[0].fields);
    } else {
        exampleEmbed.setDescription(brMessage.content);
    }
    // Send embed to BoostRequest backend THEN add the Thumbsup Icon
    const message = await (await client.channels.fetch(channelId)).send(
        exampleEmbed
    );

    if (brMessage.embeds.length == 0) {
        exampleEmbed.addField(
            "Requested By",
            `<@${brMessage.id}> ${brMessage.author.tag}`
        );
        client.channels
            .fetch(config.logChannel)
            .then((channel) => channel.send(exampleEmbed))
            .catch(console.error);
    }
    return message;
}

async function sendEmbed(
    embedUser,
    { requesterId, buyerDiscordName, message },
    { notifyBuyer, backendId }
) {
    // Make Embed post here
    const requestUser = await client.users.fetch(requesterId).catch(() => null);
    const isRealUser = requestUser && !requestUser.bot;
    const announcementEmbed = new Discord.MessageEmbed()
        .setColor("#FF0000")
        .setThumbnail(embedUser?.displayAvatarURL())
        .setTitle("An advertiser has been selected.")
        .setDescription(
            isRealUser
                ? `<@${embedUser.id}> will handle the following boost request.`
                : `<@${embedUser.id}> will handle ${buyerDiscordName}'s boost request.`
        );
    const advertiserDMEmbed = new Discord.MessageEmbed()
        .setColor("#FF0000")
        .setThumbnail(requestUser?.displayAvatarURL())
        .setTitle("You have been selected to handle a boost request.")
        .setDescription(
            isRealUser
                ? `Please message <@${requesterId}> (${requestUser.tag}).`
                : `Please message ${buyerDiscordName} (battletag).`
        )
        .setTimestamp()
        .setFooter(
            "Huokan Boosting Community",
            "https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png"
        );
    if (isRealUser) {
        advertiserDMEmbed.addField("Boost Request", message);
        announcementEmbed.addField("Boost Request", message);
    }
    try {
        await embedUser.send(advertiserDMEmbed);
    } catch (err) {
        if (err.code === 50007) {
            // Cannot send messages to this user
            announcementEmbed.setDescription(
                `<@${embedUser.id}>, I can't DM you. Please allow DMs from server members by right clicking the server and enabling "Allow direct messages from server members." in Privacy Settings.\n\n${advertiserDMEmbed.description}`
            );
        }
    }
    const backendChannel = await client.channels.fetch(backendId);
    await backendChannel.send(announcementEmbed);
    if (notifyBuyer) {
        // Make Embed post here
        const selectionBREmbed = new Discord.MessageEmbed()
            .setColor("#00FF00")
            .setTitle("Huokan Boosting Community Boost Request")
            .setThumbnail(embedUser.avatarURL())
            .addFields({
                name: "An advertiser has been found.",
                value: `<@${embedUser.id}> (${embedUser.username}#${embedUser.discriminator}) will reach out to you shortly. Anyone else that messages you regarding this boost request is not from Huokan and may attempt to scam you.`,
            })
            .setTimestamp()
            .setFooter(
                "Huokan Boosting Community",
                "https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png"
            );
        const requester = await client.users.fetch(requesterId);
        const dmChannel = requester.dmChannel || (await requester.createDM());
        await dmChannel.send(selectionBREmbed);
    }
}

function getRandomAdvertiserWeighted(advertisers) {
    const totalWeight = advertisers.reduce(
        (total, advertiser) => total + (advertiser.weight ?? 1),
        0
    );
    const chosenOffset = Math.random() * totalWeight;
    let sum = 0;
    for (const advertiser of advertisers) {
        sum += advertiser.weight ?? 1;
        if (chosenOffset < sum) {
            return advertiser;
        }
    }
    return advertisers[advertisers.length - 1];
}

async function findUser(userQuery, guildId) {
    if (/^[0-9]+$/g.test(userQuery)) {
        try {
            return await client.users.fetch(userQuery);
        } catch (err) {
            // continue
        }
    }
    try {
        if (userQuery.includes("#")) {
            const usernameUser = client.users.cache.find(
                (user) => user.tag.toLowerCase() == userQuery.toLowerCase()
            );
            if (usernameUser) {
                return usernameUser;
            }
        }
        const guild = await client.guilds.fetch(guildId);
        const nicknameUser = guild.members.cache.find(
            (member) =>
                (member.nickname ?? member.user.username).toLowerCase() ==
                userQuery
        )?.user;
        return nicknameUser;
    } catch (err) {
        console.error(err);
    }
}

let destroyed = false;
function destroy() {
    if (!destroyed) {
        destroyed = true;
        client.destroy();
        saveBoostRequests();
    }
}

process.on("SIGINT", destroy);
process.on("SIGTERM", destroy);

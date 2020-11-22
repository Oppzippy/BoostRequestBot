const config = require('./config.js');
const fs = require('fs');
const Discord = require('discord.js');
const serialize = require('serialize-javascript');
const client = new Discord.Client({
	partials: ['MESSAGE', 'USER', 'REACTION', 'GUILD_MEMBER'],
});
const reactionArray = ['👍'];
const boostRequestsBySignupMessageId = new Map();
const boostRequestTimeouts = new Map();
client.login(config.TOKEN);
let areBoostRequestsLoaded = false;
client.on('ready', () => {
	console.log('client online!');
	if (!areBoostRequestsLoaded) {
		areBoostRequestsLoaded = true;
		loadBoostRequests();
	}
});

function loadBoostRequests() {
	fs.readFile(`${__dirname}/boost-requests.js`, (err, loadedBoostRequests) => {
		if (err) {
			console.error('Failed to load boost requests', err);
			return;
		}
		try {
			const boostRequests = eval(`(${loadedBoostRequests.toString()})`);
			boostRequests.forEach((boostRequest, key) => {
				boostRequestsBySignupMessageId.set(key, boostRequest);
				addTimers(boostRequest);
			});
		}
		catch (err) {
			console.error('Failed to parse boost requests', err);
		}
	});
}

function saveBoostRequests() {
	const serialziedBoostRequests = serialize(boostRequestsBySignupMessageId, { unsafe: true, ignoreFunction: true });
	fs.writeFileSync(`${__dirname}/boost-requests.js`, serialziedBoostRequests);
}

// TODO: Edit first embed to get of 2nd embed
// TODO: Clean up embeds
// TODO: Fix bug so it doesnt respond to one person and mark all with checks
// TODO: Put every reaction.message.id into an array and check again the current one being operated on so we know if we should be using it or not


// Event Catcher when users react to messages
client.on('messageReactionAdd', async (reaction, user) => {
	try {
		if (reaction.partial) {
			await reaction.fetch();
		}
		if (user.partial) {
			await user.fetch();
		}

	}
	catch (err) {
		console.error(reaction.id, user.id, err);
		return;
	}
	console.log(`${user.username} reacted, doing checks`);
	const signupMessage = boostRequestsBySignupMessageId.get(reaction.message.id);
	const guildMember = await reaction.message.guild.members.fetch(user);
	console.log(`${signupMessage ? 'signupMessage is defined.' : 'Signup message is undefined! ' + reaction.message.id}`);

	if (signupMessage && !user.bot && reaction.emoji.name === '👍') {
		const isEliteAdvertiser = guildMember.roles.cache.some(role => role.name === 'Elite Advertiser');
		console.log(`${user.username} reacted (${isEliteAdvertiser ? '' : 'not '}elite advertiser)`);
		if (
			signupMessage.isClaimableByAdvertisers ||
			isEliteAdvertiser
		) {
			await setWinner(reaction.message, user);
		}
		else {
			signupMessage.queuedAdvertiserIds.add(user.id);
		}
	}
});

client.on('messageReactionRemove', async (reaction, user) => {
	if (reaction.partial) {
		await reaction.fetch();
	}
	const boostRequest = boostRequestsBySignupMessageId.get(reaction.message.id);
	if (boostRequest) {
		await boostRequest.queuedAdvertiserIds.delete(user.id);
	}
});

// Event Catcher when users send a message
client.on('message', async message => {
	if (message.partial) {
		await message.fetch();
	}
	if (message.author.equals(client.user)) return;
	console.log(message.content);
	const boostRequestChannel = config.BOOST_REQUEST_CHANNEL_ID.find(chan => chan.id == message.channel.id);
	// If User is not a bot AND is messsaging in BoostRequest Channel
	if (boostRequestChannel && (!message.author.bot || !boostRequestChannel.notifyBuyer)) {
		// Create embed in the Backend Channel
		if (!boostRequestChannel.useBuyerMessage) {
			if (!(await sendBuyerWaitingMessage(message))) {
				return;
			}
		}
		const signupMessage = boostRequestChannel.useBuyerMessage
			? message
			: await BREmbed(message, boostRequestChannel.backendId);
		shuffle(reactionArray);
		const reactPromises = reactionArray.map(emoji => signupMessage.react(emoji));
		await Promise.all(reactPromises);

		const buyerDiscordName = message.embeds.length >= 1
			? message.embeds[0].fields.find(field => field.name.toLowerCase().includes('battletag'))?.value
			: undefined;
		const boostRequest = {
			channelId: message.channel.id,
			requesterId: message.author.id,
			createdAt: message.createdAt,
			backendChannelId: boostRequestChannel.backendId,
			buyerDiscordName: buyerDiscordName,
			isClaimableByAdvertisers: false,
			queuedAdvertiserIds: new Set(),
			signupMessageId: signupMessage.id,
			message: message.content,
		};
		addTimers(boostRequest);
		boostRequestsBySignupMessageId.set(signupMessage.id, boostRequest);
	}
});

async function sendBuyerWaitingMessage(message) {
	const embed = new Discord.MessageEmbed()
		.setTitle('Huokan Boosting Community Boost Request')
		.setDescription(message.content)
		.setThumbnail(message.author.avatarURL())
		.setAuthor(`${message.author.username}#${message.author.discriminator}`)
		.setFooter('Huokan Boosting Community', 'https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png')
		.setTimestamp();
	try {
		const dmChannel = message.author.dmChannel ?? await message.author.createDM();
		await dmChannel.send(
			'Please wait while we find an advertiser to complete your request.',
			embed,
		);
		if (message.deletable) {
			await message.delete();
		}
	}
	catch (err) {
		if (err.code === 50007) {
			// Cannot send messages to this user
			const reply = await message.reply('I can\'t DM you! Please allow DMs from server members by right clicking the server and enabling "Allow direct messages from server members." in Privacy Settings, and then post your message again.');
			setTimeout(() => {
				message.delete().catch(() => {
					// ignore
				});
				reply.delete().catch(() => {
					// ignore
				});
			}, 30000);
		}
		else {
			console.error(err);
		}
		return false;
	}
	return true;
}

function addTimers(boostRequest) {
	boostRequestTimeouts.set(boostRequest, [
		setTimeout(async () => {
			boostRequest.isClaimableByAdvertisers = true;
			if (boostRequest.queuedAdvertiserIds.size >= 1) {
				try {
					const userId = boostRequest.queuedAdvertiserIds.values().next().value;
					const user = await client.users.fetch(userId);
					const channel = await client.channels.fetch(boostRequest.backendChannelId);
					const signupMessage = await channel.messages.fetch(boostRequest.signupMessageId);
					await setWinner(signupMessage, user);
				}
				catch (err) {
					console.error(err);
				}
			}
		}, Math.max(0, 60000 - (new Date() - boostRequest.createdAt))),
		// 1 minute
		setTimeout(() => {
			console.log('Deleting expired boost request.');
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
	boostRequestTimeouts.get(signupMessage).forEach(timeoutId => clearTimeout(timeoutId));
	const boostRequestChannel = config.BOOST_REQUEST_CHANNEL_ID.find(
		chan => chan.id == signupMessage.channelId,
	);

	try {
		const winnerName = winner.username;
		console.log(winnerName + ' won!');
		// remove reactions.
		try {
			await message.reactions.removeAll();
		}
		catch (err) {
			console.error('Failed to clear reactions: ', err);
		}
		await message.react('✅');
		await sendEmbed(winner, signupMessage, boostRequestChannel);
	}
	catch (error) {
		console.error('One of the emojis failed to react.', error);
	}
}


async function BREmbed(brMessage, channelId) {
	// Variable to eaily add hyperlink to the original message.
	const exampleEmbed = new Discord.MessageEmbed()
		.setColor('#0000FF')
		.setTitle('New Boost Request')
		.setTimestamp()
		.setFooter('Huokan Boosting Community', 'https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png');
	if (brMessage.embeds.length >= 1) {
		exampleEmbed.addFields(brMessage.embeds[0].fields);
	}
	else {
		exampleEmbed.setDescription(brMessage.content);
	}
	// Send embed to BoostRequest backend THEN add the Thumbsup Icon
	const message = await (await client.channels.fetch(channelId)).send(exampleEmbed);
	return message;
}


async function sendEmbed(embedUser, { requesterId, buyerDiscordName, message }, { notifyBuyer, backendId }) {
	// Make Embed post here
	const requestUser = await client.users.fetch(requesterId).catch(() => null);
	const isRealUser = requestUser && !requestUser.bot;
	const selectionBRBEmbed = new Discord.MessageEmbed()
		.setColor('#FF0000')
		.setThumbnail(requestUser?.displayAvatarURL())
		.setTitle('You have been selected to handle a boost request.')
		.setDescription(
			isRealUser ?
				`Please message <@${requesterId}> (${requestUser.tag}).` :
				`Please message ${buyerDiscordName} (battletag).`,
		)
		.setTimestamp()
		.setFooter('Huokan Boosting Community', 'https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png');
	if (isRealUser) {
		selectionBRBEmbed.addField('Boost Request', message);
	}
	try {
		await embedUser.send(selectionBRBEmbed);
	}
	catch (err) {
		if (err.code === 50007) {
			// Cannot send messages to this user
			const backendChannel = await client.channels.fetch(backendId);
			selectionBRBEmbed.setTitle(`${embedUser.nickname ?? embedUser.tag} has been chosen to handle a boost request.`);
			selectionBRBEmbed.setDescription(`<@${embedUser.id}>, I can't DM you. Please allow DMs from server members by right clicking the server and enabling "Allow direct messages from server members." in Privacy Settings.\n\n${selectionBRBEmbed.description}`);
			await backendChannel.send(selectionBRBEmbed);
		}
	}
	if (notifyBuyer) {
		// Make Embed post here
		const selectionBREmbed = new Discord.MessageEmbed()
			.setColor('#00FF00')
			.setTitle('Huokan Boosting Community Boost Request')
			.setThumbnail(embedUser.avatarURL())
			.addFields(
				{ name: 'An advertiser has been found.', value: `<@${embedUser.id}> (${embedUser.username}#${embedUser.discriminator}) will reach out to you shortly. Anyone else that messages you regarding this boost request is not from Huokan and may attempt to scam you.` })
			.setTimestamp()
			.setFooter('Huokan Boosting Community', 'https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png');
		const requester = await client.users.fetch(requesterId);
		const dmChannel = requester.dmChannel || await requester.createDM();
		await dmChannel.send(selectionBREmbed);
	}
}

function shuffle(array) {
	array.sort(() => Math.random() - 0.5);
}

let destroyed = false;
function destroy() {
	if (!destroyed) {
		destroyed = true;
		client.destroy();
		saveBoostRequests();
	}
}

process.on('SIGINT', destroy);
process.on('SIGTERM', destroy);

const config = require('./config.js');
const Discord = require('discord.js');
const client = new Discord.Client();
const reactionArray = ['👍'];
const boostRequestsBySignupMessageId = new Map();
client.login(config.TOKEN);
client.on('ready', () => {
	console.log('client online!');
});

// TODO: Edit first embed to get of 2nd embed
// TODO: Clean up embeds
// TODO: Fix bug so it doesnt respond to one person and mark all with checks
// TODO: Put every reaction.message.id into an array and check again the current one being operated on so we know if we should be using it or not


// Event Catcher when users react to messages
client.on('messageReactionAdd', async (reaction, user) => {
	// if user is not the bot + reaction was in backend channel + Confirm reaction is Thumbs Up
	const signupMessage = boostRequestsBySignupMessageId.get(reaction.message.id);
	const guildMember = await reaction.message.guild.members.fetch(user.id);
	console.log(user.username + ' reacted');
	if (signupMessage && !user.bot && reaction.emoji.name === '👍') {
		if (
			signupMessage.isClaimableByAdvertisers ||
			guildMember.roles.cache.some(role => role.name === 'Elite Advertiser')
		) {
			setWinner(reaction.message, user);
		}
		else {
			signupMessage.queuedAdvertiserIds.add(user.id);
		}
	}
});

client.on('messageReactionRemove', (reaction, user) => {
	const boostRequest = boostRequestsBySignupMessageId.get(reaction.message.id);
	if (boostRequest) {
		boostRequest.queuedAdvertiserIds.delete(user.id);
	}
});

// Event Catcher when users send a message
client.on('message', async message => {
	if (message.author.equals(client.user)) return;
	console.log(message.content);
	const boostRequestChannel = config.BOOST_REQUEST_CHANNEL_ID.find(chan => chan.id == message.channel.id);
	// If User is not a bot AND is messsaging in BoostRequest Channel
	if (boostRequestChannel && (!message.author.bot || !boostRequestChannel.notifyBuyer)) {
		// Create embed in the Backend Channel
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
			buyerMessageId: boostRequestChannel.useBuyerMessage ? message.id : undefined,
		};
		addTimers(boostRequest);
		boostRequestsBySignupMessageId.set(signupMessage.id, boostRequest);
		if (!boostRequestChannel.useBuyerMessage) {
			if (message.deletable) {
				message.delete();
			}
			const dmChannel = message.author.dmChannel ?? await message.author.createDM();
			const embed = new Discord.MessageEmbed()
				.setTitle('Huokan Boosting Community Boost Request')
				.setDescription(message.content)
				.setThumbnail(message.author.avatarURL())
				.setAuthor(`${message.author.username}#${message.author.discriminator}`)
				.setFooter('Huokan Boosting Community', 'https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png')
				.setTimestamp();
			dmChannel.send(
				'Please wait while we find an advertiser to complete your request.',
				embed,
			);
		}
	}
});

function addTimers(boostRequest) {
	boostRequest.timeoutIds = [
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
		}, 60000),
		setTimeout(() => {
			boostRequestsBySignupMessageId.delete(boostRequest.signupMessageId);
		}, 259200000),
	];
}

async function setWinner(message, winner) {
	const signupMessage = boostRequestsBySignupMessageId.get(message.id);
	if (!signupMessage) {
		return;
	}
	boostRequestsBySignupMessageId.delete(message.id);
	signupMessage.timeoutIds.forEach(timeoutId => clearTimeout(timeoutId));
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
		console.log(message.reactions.cache.has('✅'));
		await sendEmbed(winner, signupMessage.requesterId, {
			notifyBuyer: boostRequestChannel.notifyBuyer,
			buyerDiscordName: signupMessage.buyerDiscordName,
			backendId: boostRequestChannel.backendId,
		});
	}
	catch (error) {
		console.error('One of the emojis failed to react.', error);
	}
}


async function BREmbed(brMessage, channelId) {
	// Variable to eaily add hyperlink to the original message.
	const messagelink = brMessage.content;
	const exampleEmbed = new Discord.MessageEmbed()
		.setColor('#0000FF')
		.setTitle('New Boost Request')
		.setThumbnail(brMessage.author.displayAvatarURL())
		.setTimestamp()
		.setFooter('Huokan Boosting Community', 'https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png');
	if (brMessage.embeds.length >= 1) {
		exampleEmbed.addFields(brMessage.embeds[0].fields);
	}
	else {
		exampleEmbed.addFields({ name: brMessage.author.username, value: messagelink });
	}
	// Send embed to BoostRequest backend THEN add the Thumbsup Icon
	const message = await (await client.channels.fetch(channelId)).send(exampleEmbed);
	return message;
}


async function sendEmbed(embedUser, requesterId, { notifyBuyer, backendId, buyerDiscordName }) {
	// Make Embed post here
	const requestUser = await client.users.fetch(requesterId).catch(() => null);
	const selectionBRBEmbed = new Discord.MessageEmbed()
		.setColor('#FF0000')
		.setTitle('And the fastest clicker is...')
		.setThumbnail(embedUser.displayAvatarURL())
		.addFields({
			name: embedUser.username,
			value: requestUser && !requestUser.bot
				? `Please message <@${requesterId}>.`
				: `Please message ${buyerDiscordName} (battletag).`,
		})
		.setTimestamp()
		.setFooter('Huokan Boosting Community', 'https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png');
	await (await client.channels.fetch(backendId)).send(selectionBRBEmbed);
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
